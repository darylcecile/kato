import { parse } from "./parser.js";
import type { Document, Value, Statement, ObjectNode } from "./types.js";

export interface EvaluateOptions {
  env?: Record<string, string>;
  profile?: string;
}

/**
 * Parses Kato input and returns a plain JavaScript object.
 */
export function evaluate(input: string, options?: EvaluateOptions): Record<string, any> {
  throw new Error("not implemented");
}
