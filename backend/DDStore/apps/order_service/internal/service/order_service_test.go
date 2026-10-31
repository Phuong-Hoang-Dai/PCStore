package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Phuong-Hoang-Dai/DDStore/app/model/catalog"
	"github.com/Phuong-Hoang-Dai/DDStore/app/order_service/internal/model"
)

type reserveResult struct {
	remaining int
	err       error
}

type fakeStockRepo struct {
	results      map[string]reserveResult
	restocked    []string
	reserveCalls int
	releaseCalls int
}

func (f *fakeStockRepo) ReserveStock(productId string, quantity int, orderId int, ctx context.Context) (int, error) {
	f.reserveCalls++
	r, ok := f.results[productId]
	if !ok {
		return 10, nil
	}
	return r.remaining, r.err
}

func (f *fakeStockRepo) ReleaseStock(orderId int, ctx context.Context) ([]string, error) {
	f.releaseCalls++
	return f.restocked, nil
}

func (f *fakeStockRepo) AddStock(productId string, quantity int, ctx context.Context) (int, error) {
	return quantity, nil
}

func (f *fakeStockRepo) GetStocks(productIds []string, ctx context.Context) ([]model.Stock, error) {
	return nil, nil
}

type fakeStockCache struct {
	quantities map[string]int
	getErr     error
}

func (f *fakeStockCache) Set(productId string, quantity int, ctx context.Context) error {
	f.quantities[productId] = quantity
	return nil
}

func (f *fakeStockCache) Delete(productIds []string, ctx context.Context) error {
	for _, id := range productIds {
		delete(f.quantities, id)
	}
	return nil
}

func (f *fakeStockCache) Get(productIds []string, ctx context.Context) (map[string]int, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	found := map[string]int{}
	for _, id := range productIds {
		if q, ok := f.quantities[id]; ok {
			found[id] = q
		}
	}
	return found, nil
}

type fakeProductService struct {
	priceCalls int
}

func (f *fakeProductService) GetProductsByIds(ids []string, ctx context.Context) ([]catalog.Product, error) {
	return nil, nil
}

func (f *fakeProductService) UpdateProductState(productId string, state catalog.ProductState, ctx context.Context) error {
	return nil
}

func (f *fakeProductService) GetPriceProduct(items *[]OrderDTO, ctx context.Context) error {
	f.priceCalls++
	return nil
}

type fakeOrderRepo struct {
	restocked []string
}

func (f fakeOrderRepo) CreateOrder(data *model.Order, ctx context.Context) (int, error) {
	return data.Id, nil
}
func (f fakeOrderRepo) UpdateOrder(data model.Order, ctx context.Context) error { return nil }
func (f fakeOrderRepo) GetOrderById(id int, ctx context.Context) (model.Order, error) {
	return model.Order{Id: id, State: model.Processing}, nil
}
func (f fakeOrderRepo) GetOrders(p model.Paging, ctx context.Context) ([]model.Order, error) {
	return nil, nil
}
func (f fakeOrderRepo) GetHistoryOrders(id int, p model.Paging, ctx context.Context) ([]model.Order, error) {
	return nil, nil
}
func (f fakeOrderRepo) CancelOrder(id int, ctx context.Context) ([]string, error) {
	return f.restocked, nil
}
func (f fakeOrderRepo) CancelExpiredOrders(createdBefore time.Time, ctx context.Context) ([]int, []string, error) {
	return nil, f.restocked, nil
}

type fakeSequencerRepo struct{}

func (fakeSequencerRepo) GetOrderId(ctx context.Context) (int, error) { return 42, nil }

type fakePayment struct{}

func (fakePayment) CreatePaymentInfo(orderId int, userId int, ctx context.Context) (int, error) {
	return 0, nil
}

type orderFixture struct {
	stock    *fakeStockRepo
	cache    *fakeStockCache
	products *fakeProductService
	orders   fakeOrderRepo
}

func newOrderFixture() *orderFixture {
	return &orderFixture{
		stock:    &fakeStockRepo{results: map[string]reserveResult{}},
		cache:    &fakeStockCache{quantities: map[string]int{}},
		products: &fakeProductService{},
	}
}

func (f *orderFixture) service() OrderService {
	return NewOrderService(
		f.orders,
		f.products,
		StockService{Repos: f.stock, ProductService: f.products},
		SequencerService{Repos: fakeSequencerRepo{}},
		fakePayment{},
		f.cache,
	)
}

func order(items ...OrderDTO) []OrderDTO { return items }

func TestCreateOrderRefusesCachedShortStock(t *testing.T) {
	tests := []struct {
		name   string
		cached int
		want   error
	}{
		{"sold out", 0, model.ErrOutOfStock},
		{"not enough", 2, model.ErrNotEnoughStock},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newOrderFixture()
			f.cache.quantities["p1"] = tt.cached

			_, err := f.service().CreateOrder(1, order(OrderDTO{ProductId: "p1", Quantity: 3}), context.Background())

			if !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
			if f.products.priceCalls != 0 || f.stock.reserveCalls != 0 {
				t.Fatalf("catalog called %d times, stock %d times; want neither", f.products.priceCalls, f.stock.reserveCalls)
			}
		})
	}
}

func TestCreateOrderWithinCachedQuantityGoesOn(t *testing.T) {
	f := newOrderFixture()
	f.cache.quantities["p1"] = 3

	id, err := f.service().CreateOrder(1, order(OrderDTO{ProductId: "p1", Quantity: 3}), context.Background())

	if err != nil || id != 42 {
		t.Fatalf("CreateOrder = %d, %v; want 42, nil", id, err)
	}
	if f.stock.reserveCalls != 1 {
		t.Fatalf("reserveCalls = %d, want 1", f.stock.reserveCalls)
	}
}

func TestCreateOrderGoesOnWhenCacheFails(t *testing.T) {
	f := newOrderFixture()
	f.cache.getErr = errors.New("redis down")

	if _, err := f.service().CreateOrder(1, order(OrderDTO{ProductId: "p1", Quantity: 1}), context.Background()); err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
}

func TestCreateOrderCachesQuantityWhenNotEnough(t *testing.T) {
	f := newOrderFixture()
	f.stock.results["p1"] = reserveResult{2, model.ErrNotEnoughStock}

	_, err := f.service().CreateOrder(1, order(OrderDTO{ProductId: "p1", Quantity: 5}), context.Background())

	if !errors.Is(err, model.ErrNotEnoughStock) {
		t.Fatalf("err = %v, want ErrNotEnoughStock", err)
	}
	if q, ok := f.cache.quantities["p1"]; !ok || q != 2 {
		t.Fatalf("cached p1 = %d (%v), want 2", q, ok)
	}
}

func TestCreateOrderKeepsTheLastUnits(t *testing.T) {
	f := newOrderFixture()
	f.stock.results["p1"] = reserveResult{0, model.ErrOutOfStock}

	id, err := f.service().CreateOrder(1, order(OrderDTO{ProductId: "p1", Quantity: 3}), context.Background())

	if err != nil || id != 42 {
		t.Fatalf("CreateOrder = %d, %v; want 42, nil", id, err)
	}
	if q, ok := f.cache.quantities["p1"]; !ok || q != 0 {
		t.Fatalf("cached p1 = %d (%v), want 0", q, ok)
	}
	if f.stock.releaseCalls != 0 {
		t.Fatalf("releaseCalls = %d, want 0", f.stock.releaseCalls)
	}
}

func TestCreateOrderRollbackUncachesRestocked(t *testing.T) {
	f := newOrderFixture()
	f.stock.results["p1"] = reserveResult{0, model.ErrOutOfStock}
	f.stock.results["p2"] = reserveResult{1, model.ErrNotEnoughStock}
	f.stock.restocked = []string{"p1"}

	_, err := f.service().CreateOrder(1, order(
		OrderDTO{ProductId: "p1", Quantity: 3},
		OrderDTO{ProductId: "p2", Quantity: 2},
	), context.Background())

	if !errors.Is(err, model.ErrNotEnoughStock) {
		t.Fatalf("err = %v, want ErrNotEnoughStock", err)
	}
	if f.stock.releaseCalls != 1 {
		t.Fatalf("releaseCalls = %d, want 1", f.stock.releaseCalls)
	}
	if _, ok := f.cache.quantities["p1"]; ok {
		t.Fatal("p1 is still cached after its units were given back")
	}
	if q := f.cache.quantities["p2"]; q != 1 {
		t.Fatalf("cached p2 = %d, want 1", q)
	}
}

func TestCancelOrderUncachesRestocked(t *testing.T) {
	f := newOrderFixture()
	f.orders.restocked = []string{"p1"}
	f.cache.quantities["p1"] = 0

	if err := f.service().CancelOrder(7, context.Background()); err != nil {
		t.Fatalf("err = %v", err)
	}
	if _, ok := f.cache.quantities["p1"]; ok {
		t.Fatal("p1 is still cached after the cancel restocked it")
	}
}
