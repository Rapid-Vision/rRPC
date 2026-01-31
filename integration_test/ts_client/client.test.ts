import { describe, expect, it } from "bun:test";

import {
	CustomRPCError,
	InputRPCError,
	NotImplementedRPCError,
	ForbiddenRPCError,
	PayloadModel,
	RPCClient,
	RPCErrorException,
	TextModel,
	UnauthorizedRPCError,
	ValidationRPCError,
} from "./rpcclient";

const baseURL = "http://localhost:8080";

describe("rpcclient", () => {
	it("handles empty response", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const empty = await rpc.testEmpty();
		expect(empty).toEqual({});
	});

	it("handles no return", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const result = await rpc.testNoReturn();
		expect(result).toBeUndefined();
	});

	it("handles basic payload", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const text: TextModel = { title: null, body: "  hello  " };
		const result = await rpc.testBasic({
			text,
			flag: true,
			count: 3,
			note: "note",
		});
		expect(result.body).toBe("hello");
		expect(result.title).toBe("note");
	});

	it("handles list/map payload", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const result = await rpc.testListMap({
			texts: [
				{ title: "t1", body: "b1" },
				{ title: "t2", body: "b2" },
			],
			flags: { mode: "fast" },
		});
		expect(result.flags?.retries).toBe(2);
		expect(result.flags?.meta?.mode).toBe("fast");
		expect(result.lookup.first?.body).toBe("b1");
	});

	it("handles optional payloads", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const result = await rpc.testOptional({ text: null, flag: null });
		expect(result.enabled).toBe(false);
	});

	it("maps validation error", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		await expect(
			rpc.testValidationError({ text: { title: null, body: "" } })
		).rejects.toBeInstanceOf(ValidationRPCError);
	});

	it("maps input error", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		await expect(
			rpc.testBasic({ text: "bad" as unknown as TextModel, flag: true, count: 1 })
		).rejects.toBeInstanceOf(InputRPCError);
	});

	it("maps unauthorized error", async () => {
		const rpc = new RPCClient(baseURL);
		await expect(rpc.testUnauthorizedError()).rejects.toBeInstanceOf(
			UnauthorizedRPCError
		);
	});

	it("maps forbidden error", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		await expect(rpc.testForbiddenError()).rejects.toBeInstanceOf(
			ForbiddenRPCError
		);
	});

	it("maps not implemented error", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		await expect(rpc.testNotImplementedError()).rejects.toBeInstanceOf(
			NotImplementedRPCError
		);
	});

	it("maps custom error", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		await expect(rpc.testCustomError()).rejects.toBeInstanceOf(CustomRPCError);
	});

	it("handles map return", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const result = await rpc.testMapReturn();
		expect(result.a?.body).toBe("mapped");
	});

	it("handles json payload", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const result = await rpc.testJson({ data: { count: 2, tags: ["a", "b"] } });
		expect(result.count).toBe(2);
	});

	it("handles raw payload", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const result = await rpc.testRaw({ payload: { ok: true } });
		expect(result.ok).toBe(true);
	});

	it("handles mixed payload", async () => {
		const rpc = new RPCClient(baseURL, {
			headers: { Authorization: "Bearer test_token" },
		});
		const payload: PayloadModel = { data: { value: "x" }, raw_data: { id: 1 } };
		const result = await rpc.testMixedPayload({ payload });
		expect(result.data.value).toBe("x");
		expect(result.raw_data.id).toBe(1);
	});

	it("normalizes base url and prefix", async () => {
		const rpc = new RPCClient("localhost:8080/", {
			prefix: "rpc",
			headers: { Authorization: "Bearer test_token" },
		});
		const empty = await rpc.testEmpty();
		expect(empty).toEqual({});
	});

	it("maps non-json http error to custom", async () => {
		const rpc = new RPCClient(baseURL, {
			fetchFn: async () => new Response("boom", { status: 500 }),
		});
		try {
			await rpc.testEmpty();
			throw new Error("expected request to fail");
		} catch (err) {
			expect(err).toBeInstanceOf(RPCErrorException);
			const rpcErr = err as RPCErrorException;
			expect(rpcErr.error.type).toBe("custom");
			expect(rpcErr.error.message).toBe("rpc error: status 500");
		}
	});
});
