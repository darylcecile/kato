import { describe, it, expect } from "vitest";
import { parse } from "../src/index.js";
import type {
  Document,
  KeyValue,
  NamedBlock,
  NamedArray,
  ArrayOp,
  ScalarNode,
  UnitLiteral,
  FunctionCall,
  ObjectNode,
  ArrayNode,
  Directive,
} from "../src/index.js";

describe("Kato Parser", () => {
  // =======================================================================
  // Basic Object Parsing
  // =======================================================================
  describe("basic objects", () => {
    it("parses a simple named block with key-value pairs", () => {
      const input = `server {
  port: 3000
  host: "localhost"
}`;
      const doc = parse(input);
      expect(doc.type).toBe("Document");
      expect(doc.body).toHaveLength(1);

      const block = doc.body[0] as NamedBlock;
      expect(block.type).toBe("NamedBlock");
      expect(block.name).toBe("server");
      expect(block.body.members).toHaveLength(2);

      const port = block.body.members[0] as KeyValue;
      expect(port.key).toBe("port");
      expect((port.value as ScalarNode).value).toBe(3000);

      const host = block.body.members[1] as KeyValue;
      expect(host.key).toBe("host");
      expect((host.value as ScalarNode).value).toBe("localhost");
    });

    it("parses an empty object", () => {
      const doc = parse(`empty {}`);
      const block = doc.body[0] as NamedBlock;
      expect(block.body.members).toHaveLength(0);
    });

    it("parses top-level key-value pairs", () => {
      const input = `name: "Daryl"\nage: 25`;
      const doc = parse(input);
      expect(doc.body).toHaveLength(2);
      expect((doc.body[0] as KeyValue).key).toBe("name");
      expect((doc.body[1] as KeyValue).key).toBe("age");
    });
  });

  // =======================================================================
  // Nested Objects
  // =======================================================================
  describe("nested objects", () => {
    it("parses deeply nested objects", () => {
      const input = `database {
  pool {
    min: 2
    max: 10
  }
}`;
      const doc = parse(input);
      const db = doc.body[0] as NamedBlock;
      expect(db.name).toBe("database");

      const pool = db.body.members[0] as NamedBlock;
      expect(pool.type).toBe("NamedBlock");
      expect(pool.name).toBe("pool");
      expect(pool.body.members).toHaveLength(2);

      const min = pool.body.members[0] as KeyValue;
      expect(min.key).toBe("min");
      expect((min.value as ScalarNode).value).toBe(2);
    });
  });

  // =======================================================================
  // Inline Arrays
  // =======================================================================
  describe("inline arrays", () => {
    it("parses comma-separated inline arrays", () => {
      const input = `regions: ["eu-west-1", "us-east-1"]`;
      const doc = parse(input);
      const kv = doc.body[0] as KeyValue;
      expect(kv.key).toBe("regions");

      const arr = kv.value as ArrayNode;
      expect(arr.type).toBe("Array");
      expect(arr.elements).toHaveLength(2);
      expect((arr.elements[0] as ScalarNode).value).toBe("eu-west-1");
      expect((arr.elements[1] as ScalarNode).value).toBe("us-east-1");
    });

    it("allows trailing commas in inline arrays", () => {
      const input = `items: [1, 2, 3,]`;
      const doc = parse(input);
      const arr = (doc.body[0] as KeyValue).value as ArrayNode;
      expect(arr.elements).toHaveLength(3);
    });
  });

  // =======================================================================
  // Block Arrays
  // =======================================================================
  describe("block arrays", () => {
    it("parses newline-separated block arrays", () => {
      const input = `plugins [
  "auth"
  "payments"
  "search"
]`;
      const doc = parse(input);
      const na = doc.body[0] as NamedArray;
      expect(na.type).toBe("NamedArray");
      expect(na.name).toBe("plugins");
      expect(na.elements.elements).toHaveLength(3);
      expect((na.elements.elements[0] as ScalarNode).value).toBe("auth");
    });

    it("parses mixed comma and newline separators", () => {
      const input = `tags [
  "a", "b"
  "c"
]`;
      const doc = parse(input);
      const na = doc.body[0] as NamedArray;
      expect(na.elements.elements).toHaveLength(3);
    });
  });

  // =======================================================================
  // Scalar Types
  // =======================================================================
  describe("scalars", () => {
    it("parses quoted strings", () => {
      const doc = parse(`name: "hello world"`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("string");
      expect(v.value).toBe("hello world");
    });

    it("parses integers", () => {
      const doc = parse(`count: 42`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("integer");
      expect(v.value).toBe(42);
    });

    it("parses negative integers", () => {
      const doc = parse(`offset: -10`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.value).toBe(-10);
    });

    it("parses floats", () => {
      const doc = parse(`ratio: 3.14`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("float");
      expect(v.value).toBe(3.14);
    });

    it("parses floats with exponent", () => {
      const doc = parse(`big: 1.5e10`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("float");
      expect(v.value).toBe(1.5e10);
    });

    it("parses boolean true", () => {
      const doc = parse(`enabled: true`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("boolean");
      expect(v.value).toBe(true);
    });

    it("parses boolean false", () => {
      const doc = parse(`disabled: false`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("boolean");
      expect(v.value).toBe(false);
    });

    it("parses null", () => {
      const doc = parse(`nothing: null`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("null");
      expect(v.value).toBe(null);
    });

    it("parses on/off literals", () => {
      const doc = parse(`feature: on`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("on_off");
      expect(v.value).toBe(true);
    });

    it("parses off literal", () => {
      const doc = parse(`debug: off`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("on_off");
      expect(v.value).toBe(false);
    });

    it("parses bare words as strings", () => {
      const doc = parse(`env: production`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("bare_word");
      expect(v.value).toBe("production");
    });

    it("parses hex integers", () => {
      const doc = parse(`color: 0xFF00AA`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.scalarType).toBe("integer");
      expect(v.value).toBe(0xFF00AA);
    });
  });

  // =======================================================================
  // Unit Literals
  // =======================================================================
  describe("unit literals", () => {
    it("parses time units", () => {
      const doc = parse(`timeout: 30s`);
      const v = (doc.body[0] as KeyValue).value as UnitLiteral;
      expect(v.type).toBe("UnitLiteral");
      expect(v.value).toBe(30);
      expect(v.unit).toBe("s");
    });

    it("parses memory units", () => {
      const doc = parse(`memory: 512MiB`);
      const v = (doc.body[0] as KeyValue).value as UnitLiteral;
      expect(v.value).toBe(512);
      expect(v.unit).toBe("MiB");
    });

    it("parses rate units with slash", () => {
      const doc = parse(`rateLimit: 100/min`);
      const v = (doc.body[0] as KeyValue).value as UnitLiteral;
      expect(v.value).toBe(100);
      expect(v.unit).toBe("/min");
    });

    it("parses float unit literals", () => {
      const doc = parse(`delay: 1.5s`);
      const v = (doc.body[0] as KeyValue).value as UnitLiteral;
      expect(v.value).toBe(1.5);
      expect(v.unit).toBe("s");
    });
  });

  // =======================================================================
  // Comments
  // =======================================================================
  describe("comments", () => {
    it("ignores line comments", () => {
      const input = `# this is a comment\nport: 3000`;
      const doc = parse(input);
      expect(doc.body).toHaveLength(1);
      expect((doc.body[0] as KeyValue).key).toBe("port");
    });

    it("ignores inline comments", () => {
      const input = `port: 3000  # the port number`;
      const doc = parse(input);
      const kv = doc.body[0] as KeyValue;
      expect(kv.key).toBe("port");
      expect((kv.value as ScalarNode).value).toBe(3000);
    });

    it("preserves doc comments in AST", () => {
      const input = `/// This documents the timeout\ntimeout: 30s`;
      const doc = parse(input);
      const kv = doc.body[0] as KeyValue;
      expect(kv.leadingComments).toBeDefined();
      expect(kv.leadingComments![0].type).toBe("DocComment");
      expect(kv.leadingComments![0].text).toBe("This documents the timeout");
    });

    it("preserves multi-line doc comments", () => {
      const input = `/// Line one\n/// Line two\nkey: "val"`;
      const doc = parse(input);
      const kv = doc.body[0] as KeyValue;
      expect(kv.leadingComments).toHaveLength(2);
    });
  });

  // =======================================================================
  // Directives
  // =======================================================================
  describe("directives", () => {
    it("parses @version directive", () => {
      const doc = parse(`@version 1`);
      expect(doc.directives).toHaveLength(1);
      expect(doc.directives[0].name).toBe("version");
      expect(doc.directives[0].value).toBe(1);
    });

    it("parses @schema directive", () => {
      const doc = parse(`@schema "https://example.com/schema.kato"`);
      expect(doc.directives[0].name).toBe("schema");
      expect(doc.directives[0].value).toBe("https://example.com/schema.kato");
    });

    it("parses @include directive", () => {
      const doc = parse(`@include "./shared.kato"`);
      expect(doc.directives[0].name).toBe("include");
      expect(doc.directives[0].value).toBe("./shared.kato");
    });

    it("parses @include with alias", () => {
      const doc = parse(`@include "./secrets.kato" as secrets`);
      expect(doc.directives[0].name).toBe("include");
    });

    it("parses @profile directive with body", () => {
      const input = `@profile production {
  server {
    workers: 8
  }
}`;
      const doc = parse(input);
      expect(doc.directives).toHaveLength(1);
      expect(doc.directives[0].name).toBe("profile");
      expect(doc.directives[0].profile).toBe("production");
      expect(doc.directives[0].body).toBeDefined();
      expect(doc.directives[0].body!.type).toBe("Object");
    });
  });

  // =======================================================================
  // Function Calls (env, ref)
  // =======================================================================
  describe("function calls", () => {
    it("parses env() with single argument", () => {
      const doc = parse(`url: env("DATABASE_URL")`);
      const v = (doc.body[0] as KeyValue).value as FunctionCall;
      expect(v.type).toBe("FunctionCall");
      expect(v.name).toBe("env");
      expect((v.args[0] as ScalarNode).value).toBe("DATABASE_URL");
    });

    it("parses env() with fallback option", () => {
      const doc = parse(`url: env("DB_URL", fallback: null)`);
      const v = (doc.body[0] as KeyValue).value as FunctionCall;
      expect(v.name).toBe("env");
      expect(v.options).toBeDefined();
      expect((v.options!["fallback"] as ScalarNode).value).toBe(null);
    });

    it("parses ref() call", () => {
      const doc = parse(`region: ref("shared.region")`);
      const v = (doc.body[0] as KeyValue).value as FunctionCall;
      expect(v.type).toBe("FunctionCall");
      expect(v.name).toBe("ref");
      expect((v.args[0] as ScalarNode).value).toBe("shared.region");
    });
  });

  // =======================================================================
  // Profiles
  // =======================================================================
  describe("profiles", () => {
    it("parses profile with overrides", () => {
      const input = `server {
  port: 3000
  workers: 2
}

@profile production {
  server {
    workers: 8
  }
}`;
      const doc = parse(input);
      expect(doc.body).toHaveLength(1);
      expect(doc.directives).toHaveLength(1);
      expect(doc.directives[0].name).toBe("profile");
      expect(doc.directives[0].profile).toBe("production");
    });
  });

  // =======================================================================
  // Array Operations
  // =======================================================================
  describe("array operations", () => {
    it("parses += array operation", () => {
      const doc = parse(`plugins += ["analytics"]`);
      const op = doc.body[0] as ArrayOp;
      expect(op.type).toBe("ArrayOp");
      expect(op.name).toBe("plugins");
      expect(op.operator).toBe("+=");
      expect(op.value.elements).toHaveLength(1);
    });

    it("parses -= array operation", () => {
      const doc = parse(`plugins -= ["debug-toolbar"]`);
      const op = doc.body[0] as ArrayOp;
      expect(op.operator).toBe("-=");
      expect((op.value.elements[0] as ScalarNode).value).toBe("debug-toolbar");
    });
  });

  // =======================================================================
  // Edge Cases
  // =======================================================================
  describe("edge cases", () => {
    it("handles trailing commas in objects", () => {
      const input = `server {\n  port: 3000,\n  host: "localhost",\n}`;
      const doc = parse(input);
      const block = doc.body[0] as NamedBlock;
      expect(block.body.members).toHaveLength(2);
    });

    it("handles empty arrays", () => {
      const doc = parse(`items: []`);
      const arr = (doc.body[0] as KeyValue).value as ArrayNode;
      expect(arr.elements).toHaveLength(0);
    });

    it("handles newline-separated array elements", () => {
      const input = `items: [\n  1\n  2\n  3\n]`;
      const doc = parse(input);
      const arr = (doc.body[0] as KeyValue).value as ArrayNode;
      expect(arr.elements).toHaveLength(3);
    });

    it("handles multiple top-level blocks", () => {
      const input = `server {\n  port: 3000\n}\n\ndatabase {\n  host: "db"\n}`;
      const doc = parse(input);
      expect(doc.body).toHaveLength(2);
    });

    it("handles string escape sequences", () => {
      const doc = parse(`msg: "hello\\nworld"`);
      const v = (doc.body[0] as KeyValue).value as ScalarNode;
      expect(v.value).toBe("hello\nworld");
    });
  });

  // =======================================================================
  // Error Cases
  // =======================================================================
  describe("error cases", () => {
    it("throws on unclosed braces", () => {
      expect(() => parse(`server {`)).toThrow();
    });

    it("throws on unclosed brackets", () => {
      expect(() => parse(`items: [1, 2`)).toThrow();
    });

    it("throws on unclosed strings", () => {
      expect(() => parse(`name: "unterminated`)).toThrow();
    });

    it("throws on invalid directive", () => {
      expect(() => parse(`@ `)).toThrow();
    });

    it("throws on unexpected tokens", () => {
      expect(() => parse(`}: invalid`)).toThrow();
    });

    it("throws on duplicate colon", () => {
      expect(() => parse(`key:: value`)).toThrow();
    });
  });
});
