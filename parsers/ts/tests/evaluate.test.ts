import { describe, it, expect } from "vitest";
import { evaluate } from "../src/evaluate.js";

describe("evaluate()", () => {
  describe("basic objects", () => {
    it("evaluates a simple named block", () => {
      expect(evaluate(`server { port: 3000 }`)).toEqual({ server: { port: 3000 } });
    });

    it("evaluates nested blocks", () => {
      expect(evaluate(`db { pool { min: 2 } }`)).toEqual({ db: { pool: { min: 2 } } });
    });
  });

  describe("arrays", () => {
    it("evaluates inline arrays", () => {
      expect(evaluate(`regions: ["a", "b"]`)).toEqual({ regions: ["a", "b"] });
    });
  });

  describe("scalars", () => {
    it("evaluates strings", () => {
      expect(evaluate(`name: "hello"`)).toEqual({ name: "hello" });
    });

    it("evaluates integers", () => {
      expect(evaluate(`count: 42`)).toEqual({ count: 42 });
    });

    it("evaluates floats", () => {
      expect(evaluate(`ratio: 3.14`)).toEqual({ ratio: 3.14 });
    });

    it("evaluates booleans", () => {
      expect(evaluate(`active: true`)).toEqual({ active: true });
    });

    it("evaluates null", () => {
      expect(evaluate(`value: null`)).toEqual({ value: null });
    });

    it("evaluates on as true", () => {
      expect(evaluate(`enabled: on`)).toEqual({ enabled: true });
    });

    it("evaluates off as false", () => {
      expect(evaluate(`disabled: off`)).toEqual({ disabled: false });
    });

    it("evaluates bare words", () => {
      expect(evaluate(`mode: production`)).toEqual({ mode: "production" });
    });
  });

  describe("unit literals", () => {
    it("evaluates unit literals to $unit object", () => {
      expect(evaluate(`timeout: 30s`)).toEqual({ timeout: { $unit: "s", value: 30 } });
    });
  });

  describe("env()", () => {
    it("resolves env() calls with provided env option", () => {
      const result = evaluate(`url: env("DB_URL")`, { env: { DB_URL: "postgres://localhost/db" } });
      expect(result).toEqual({ url: "postgres://localhost/db" });
    });
  });

  describe("ref()", () => {
    it("resolves ref() to the referenced value", () => {
      const input = `
shared { region: "us-east-1" }
deploy { region: ref("shared.region") }
`;
      expect(evaluate(input)).toEqual({
        shared: { region: "us-east-1" },
        deploy: { region: "us-east-1" },
      });
    });
  });

  describe("profiles", () => {
    it("merges profile overrides when profile option is set", () => {
      const input = `
server { port: 3000 }
#[profile "production"]
server { port: 8080 }
`;
      const result = evaluate(input, { profile: "production" });
      expect(result).toEqual({ server: { port: 8080 } });
    });
  });

  describe("top-level keys", () => {
    it("handles multiple top-level key-value pairs", () => {
      const input = `
host: "localhost"
port: 5432
`;
      expect(evaluate(input)).toEqual({ host: "localhost", port: 5432 });
    });
  });
});
