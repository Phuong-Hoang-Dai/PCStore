import { useCallback, useEffect, useRef, useState } from "react";
import type { StatusResponse, TargetResult } from "../server/types";

export const POLL_MS = 10_000;
export const HISTORY_LENGTH = 30;

export type Check = {
  at: number;
  up: boolean;
  reachable: boolean;
  latencyMs: number;
  error?: string;
};

export type TargetState = {
  name: string;
  group: string;
  optional: boolean;
  url: string;
  latest: Check;
  // Oldest first, at most HISTORY_LENGTH.
  history: Check[];
  // When the target last switched between up and down.
  since: number;
  // When this page first saw the target; since === firstSeen means no change was seen yet.
  firstSeen: number;
};

export type Display = "up" | "down" | "off";

export function displayOf(target: { optional: boolean }, check: Check): Display {
  if (check.up) {
    return "up";
  }
  // Nothing answering an optional target just means its compose profile isn't started.
  return target.optional && !check.reachable ? "off" : "down";
}

// useStatus polls /api/status and keeps the last HISTORY_LENGTH results of every target.
export function useStatus() {
  const [targets, setTargets] = useState<TargetState[]>([]);
  const [checkedAt, setCheckedAt] = useState<number | null>(null);
  const [apiError, setApiError] = useState("");
  const [paused, setPaused] = useState(false);
  const [checking, setChecking] = useState(false);
  // Counts finished checks; each one schedules the next.
  const [finished, setFinished] = useState(0);
  const inFlight = useRef(false);

  const check = useCallback(async () => {
    if (inFlight.current) {
      return;
    }
    inFlight.current = true;
    setChecking(true);
    try {
      const res = await fetch("/api/status", { cache: "no-store" });
      if (!res.ok) {
        throw new Error(`the status API answered ${res.status}`);
      }
      const data: StatusResponse = await res.json();
      const at = Date.parse(data.checkedAt);
      setTargets((prev) => merge(prev, data.results, at));
      setCheckedAt(at);
      setApiError("");
    } catch (err) {
      setApiError(err instanceof Error ? err.message : String(err));
    } finally {
      inFlight.current = false;
      setChecking(false);
      setFinished((n) => n + 1);
    }
  }, []);

  // The first check runs right away; each next one POLL_MS after the previous finished, so
  // checks never overlap.
  useEffect(() => {
    if (paused) {
      return;
    }
    const id = window.setTimeout(check, finished === 0 ? 0 : POLL_MS);
    return () => window.clearTimeout(id);
  }, [finished, paused, check]);

  const togglePaused = useCallback(() => {
    if (paused) {
      check();
    }
    setPaused(!paused);
  }, [paused, check]);

  return { targets, checkedAt, apiError, paused, checking, check, togglePaused };
}

function merge(prev: TargetState[], results: TargetResult[], at: number): TargetState[] {
  const byName = new Map(prev.map((t) => [t.name, t]));
  return results.map((r) => {
    const latest: Check = {
      at,
      up: r.up,
      reachable: r.reachable,
      latencyMs: r.latencyMs,
      error: r.error,
    };
    const old = byName.get(r.name);
    return {
      name: r.name,
      group: r.group,
      optional: r.optional,
      url: r.url,
      latest,
      history: [...(old?.history ?? []), latest].slice(-HISTORY_LENGTH),
      since: old && old.latest.up === r.up ? old.since : at,
      firstSeen: old?.firstSeen ?? at,
    };
  });
}

// useNow re-renders every intervalMs, for "checked 4s ago" style labels.
export function useNow(intervalMs = 1000) {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    const id = window.setInterval(() => setNow(Date.now()), intervalMs);
    return () => window.clearInterval(id);
  }, [intervalMs]);
  return now;
}
