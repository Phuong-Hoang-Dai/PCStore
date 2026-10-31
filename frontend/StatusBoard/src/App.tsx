import { useEffect, useState } from "react";
import type { ButtonHTMLAttributes } from "react";
import { HISTORY_LENGTH, POLL_MS, displayOf, useNow, useStatus } from "./useStatus";
import type { Check, Display, TargetState } from "./useStatus";

const DISPLAY: Record<Display, { label: string; icon: string; fill: string }> = {
  up: { label: "Up", icon: "✓", fill: "bg-good" },
  down: { label: "Down", icon: "✕", fill: "bg-critical" },
  off: { label: "Off", icon: "–", fill: "bg-muted" },
};

type Hover = { target: TargetState; check: Check; rect: DOMRect } | null;

const App = () => {
  const { targets, checkedAt, apiError, paused, checking, check, togglePaused } = useStatus();
  const now = useNow();
  const [hover, setHover] = useState<Hover>(null);

  const requiredDown = targets.filter((t) => !t.optional && !t.latest.up);
  const count = (display: Display) =>
    targets.filter((t) => displayOf(t, t.latest) === display).length;
  const slowest = targets
    .filter((t) => t.latest.up)
    .reduce<TargetState | null>((a, t) => (!a || t.latest.latencyMs > a.latest.latencyMs ? t : a), null);
  const groups = [...new Set(targets.map((t) => t.group))].map((name) => ({
    name,
    targets: targets.filter((t) => t.group === name),
  }));

  useEffect(() => {
    if (checkedAt === null) {
      document.title = "DDStore status";
    } else if (requiredDown.length > 0) {
      document.title = `✕ ${requiredDown.length} down · DDStore status`;
    } else {
      document.title = "✓ DDStore status";
    }
  }, [checkedAt, requiredDown.length]);

  let checkedLine = checking ? "Running the first check…" : "No results yet";
  if (checkedAt !== null) {
    const ago = formatDuration(now - checkedAt);
    checkedLine = paused
      ? `Paused · last checked ${ago} ago`
      : `Checked ${ago} ago · every ${POLL_MS / 1000}s`;
  }

  return (
    <main className="mx-auto max-w-6xl px-4 py-8">
      <header className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">DDStore status</h1>
          <p className="mt-1 text-sm text-ink-2">{checkedLine}</p>
        </div>
        <div className="flex gap-2">
          <Button onClick={check} disabled={checking}>
            {checking ? "Checking…" : "Check now"}
          </Button>
          <Button onClick={togglePaused}>{paused ? "Resume" : "Pause"}</Button>
        </div>
      </header>

      <Overall hasData={checkedAt !== null} requiredDown={requiredDown} apiError={apiError} />

      {targets.length > 0 && (
        <div className="mt-4 grid grid-cols-2 gap-4 md:grid-cols-4">
          <StatTile label="Up" value={count("up")} caption={`of ${targets.length} targets`} />
          <StatTile
            label="Down"
            value={count("down")}
            caption={`${requiredDown.length} of them required`}
          />
          <StatTile label="Not running" value={count("off")} caption="optional, profile not started" />
          <StatTile
            label="Slowest response"
            value={slowest ? `${slowest.latest.latencyMs} ms` : "–"}
            caption={slowest?.name}
          />
        </div>
      )}

      {groups.map((g) => (
        <GroupTable key={g.name} name={g.name} targets={g.targets} now={now} onHover={setHover} />
      ))}

      <Tooltip hover={hover} />
    </main>
  );
};

const Overall = ({
  hasData,
  requiredDown,
  apiError,
}: {
  hasData: boolean;
  requiredDown: TargetState[];
  apiError: string;
}) => {
  let display: Display | null = "up";
  let title = "All required targets are up";
  let detail = "Optional targets don't count: they only run with their compose profile.";
  if (apiError) {
    display = "down";
    title = "Can't reach the status API";
    detail = `${apiError}. Is npm run dev still running?${hasData ? " Showing the last results." : ""}`;
  } else if (!hasData) {
    display = null;
    title = "Checking…";
    detail = "Waiting for the first results.";
  } else if (requiredDown.length > 0) {
    display = "down";
    title = `${requiredDown.length} required ${requiredDown.length === 1 ? "target is" : "targets are"} down`;
    detail = requiredDown.map((t) => t.name).join(", ");
  }

  return (
    <section className="mt-6 flex items-center gap-4 rounded-xl bg-surface p-5 ring-1 ring-edge">
      <StatusIcon display={display} size="size-9 text-lg" />
      <div>
        <div className="text-lg font-semibold">{title}</div>
        <div className="text-sm text-ink-2">{detail}</div>
      </div>
    </section>
  );
};

const StatTile = ({ label, value, caption }: { label: string; value: number | string; caption?: string }) => (
  <div className="rounded-xl bg-surface p-4 ring-1 ring-edge">
    <div className="text-sm text-ink-2">{label}</div>
    <div className="mt-1 text-3xl font-semibold">{value}</div>
    <div className="mt-1 truncate text-sm text-ink-2">{caption ?? " "}</div>
  </div>
);

const GroupTable = ({
  name,
  targets,
  now,
  onHover,
}: {
  name: string;
  targets: TargetState[];
  now: number;
  onHover: (hover: Hover) => void;
}) => (
  <section className="mt-8">
    <h2 className="mb-2 text-sm font-semibold uppercase tracking-wide text-ink-2">{name}</h2>
    <div className="overflow-x-auto rounded-xl bg-surface ring-1 ring-edge">
      <table className="w-full text-sm">
        <thead className="text-left text-ink-2">
          <tr className="border-b border-line">
            <th className="px-4 py-2 font-medium">Status</th>
            <th className="px-4 py-2 font-medium">Target</th>
            <th className="hidden px-4 py-2 font-medium md:table-cell">Last {HISTORY_LENGTH} checks</th>
            <th className="px-4 py-2 text-right font-medium">Latency</th>
            <th className="px-4 py-2 text-right font-medium" title="How long the target has been up or down">
              For
            </th>
          </tr>
        </thead>
        <tbody>
          {targets.map((t) => (
            <TargetRow key={t.name} target={t} now={now} onHover={onHover} />
          ))}
        </tbody>
      </table>
    </div>
  </section>
);

const TargetRow = ({
  target,
  now,
  onHover,
}: {
  target: TargetState;
  now: number;
  onHover: (hover: Hover) => void;
}) => {
  const display = displayOf(target, target.latest);
  // Without a change since the page opened, it has been in this state at least this long.
  const unchanged = target.since === target.firstSeen;

  return (
    <tr className="border-b border-line align-top last:border-0">
      <td className="w-28 px-4 py-3">
        <StatusBadge display={display} />
      </td>
      <td className="px-4 py-3">
        <div className="font-medium">
          {target.name}
          {target.optional && <span className="ml-2 text-xs font-normal text-ink-2">optional</span>}
        </div>
        <div className="font-mono text-xs break-all text-ink-2">{target.url}</div>
        {target.latest.error && <div className="mt-1 text-ink">{target.latest.error}</div>}
      </td>
      <td className="hidden px-4 py-3 md:table-cell">
        <HistoryStrip target={target} onHover={onHover} />
      </td>
      <td className="px-4 py-3 text-right whitespace-nowrap tabular-nums">
        {target.latest.up ? `${target.latest.latencyMs} ms` : "–"}
      </td>
      <td
        className="px-4 py-3 text-right whitespace-nowrap text-ink-2 tabular-nums"
        title={unchanged ? "No change since this page opened" : undefined}
      >
        {unchanged ? "≥ " : ""}
        {formatDuration(now - target.since)}
      </td>
    </tr>
  );
};

// One thin cell per check, oldest on the left; empty slots keep the strip a fixed width.
const HistoryStrip = ({ target, onHover }: { target: TargetState; onHover: (hover: Hover) => void }) => {
  const up = target.history.filter((c) => c.up).length;
  return (
    <div
      role="img"
      aria-label={`Last ${target.history.length} checks: ${up} up, ${target.history.length - up} not up`}
      className="flex h-5 gap-0.5"
    >
      {Array.from({ length: HISTORY_LENGTH - target.history.length }, (_, i) => (
        <span key={`empty-${i}`} className="w-1.5 rounded-xs bg-line" />
      ))}
      {target.history.map((check) => (
        <span
          key={check.at}
          className={`w-1.5 rounded-xs ${DISPLAY[displayOf(target, check)].fill}`}
          onMouseEnter={(e) => onHover({ target, check, rect: e.currentTarget.getBoundingClientRect() })}
          onMouseLeave={() => onHover(null)}
        />
      ))}
    </div>
  );
};

const Tooltip = ({ hover }: { hover: Hover }) => {
  if (!hover) {
    return null;
  }
  const { target, check, rect } = hover;
  return (
    <div
      role="tooltip"
      className="pointer-events-none fixed z-10 max-w-72 -translate-x-1/2 -translate-y-full rounded-lg bg-surface px-3 py-2 text-sm shadow-lg ring-1 ring-edge"
      style={{ left: rect.left + rect.width / 2, top: rect.top - 6 }}
    >
      <div className="font-medium">
        {target.name} · {new Date(check.at).toLocaleTimeString()}
      </div>
      <div className="text-ink-2">
        {DISPLAY[displayOf(target, check)].label} · {check.latencyMs} ms
      </div>
      {check.error && <div className="mt-1 text-ink-2">{check.error}</div>}
    </div>
  );
};

const StatusBadge = ({ display }: { display: Display }) => (
  <span className="inline-flex items-center gap-2 font-medium">
    <StatusIcon display={display} />
    {DISPLAY[display].label}
  </span>
);

// null is "not checked yet".
const StatusIcon = ({ display, size = "size-5 text-xs" }: { display: Display | null; size?: string }) => {
  if (!display) {
    return <span aria-hidden className={`${size} shrink-0 animate-pulse rounded-full bg-line`} />;
  }
  const { icon, fill } = DISPLAY[display];
  return (
    <span aria-hidden className={`grid ${size} shrink-0 place-items-center rounded-full font-bold text-white ${fill}`}>
      {icon}
    </span>
  );
};

const Button = (props: ButtonHTMLAttributes<HTMLButtonElement>) => (
  <button
    type="button"
    {...props}
    className="rounded-lg bg-surface px-3 py-1.5 text-sm font-medium ring-1 ring-edge hover:bg-line disabled:cursor-not-allowed disabled:opacity-60"
  />
);

function formatDuration(ms: number) {
  const s = Math.max(0, Math.round(ms / 1000));
  if (s < 60) {
    return `${s}s`;
  }
  const m = Math.floor(s / 60);
  if (m < 60) {
    return `${m}m`;
  }
  return `${Math.floor(m / 60)}h ${m % 60}m`;
}

export default App;
