import { describe, expect, it } from "bun:test";
import { ZodError } from "zod";

import {
	RPCClient,
	TextModelSchema,
	TestOptionalParamsSchema,
} from "./rpcclientzod";
import type { TextModel } from "./rpcclientzod";

const baseURL = "http://localhost:8080";

describe("rpcclientzod", () => {
	it("validates inputs before sending", async () => {
		const rpc = new RPCClient(baseURL, {
			bearerToken: "test_token",
			fetchFn: async () => {
				throw new Error("fetch should not be called");
			},
		});

		await expect(
			rpc.testBasic({
				text: { title: "missing-body" } as unknown as TextModel,
				flag: true,
				count: 1,
			})
		).rejects.toBeInstanceOf(ZodError);
	});

	it("accepts optional nullable fields", () => {
		expect(() =>
			TestOptionalParamsSchema.parse({
				text: null,
				flag: null,
			})
		).not.toThrow();
		expect(() =>
			TextModelSchema.parse({
				title: null,
				body: "ok",
			})
		).not.toThrow();
	});
});
