// THIS CODE IS GENERATED

import { ERROR_EXCEPTIONS, RPCErrorException } from "./errors";
import type { RPCError } from "./errors";
import type {
	EmptyModel,
	TextModel,
	FlagsModel,
	NestedModel,
	PayloadModel,
	TestEmptyResult,
	TestBasicParams,
	TestBasicResult,
	TestListMapParams,
	TestListMapResult,
	TestOptionalParams,
	TestOptionalResult,
	TestValidationErrorParams,
	TestValidationErrorResult,
	TestUnauthorizedErrorResult,
	TestForbiddenErrorResult,
	TestNotImplementedErrorResult,
	TestCustomErrorResult,
	TestMapReturnResult,
	TestJsonParams,
	TestJsonResult,
	TestRawParams,
	TestRawResult,
	TestMixedPayloadParams,
	TestMixedPayloadResult,
} from "./models";
import {
	TestBasicParamsSchema,
	TestListMapParamsSchema,
	TestOptionalParamsSchema,
	TestValidationErrorParamsSchema,
	TestJsonParamsSchema,
	TestRawParamsSchema,
	TestMixedPayloadParamsSchema,
} from "./models";

export type FetchResponse = {
	ok: boolean;
	status: number;
	json(): Promise<unknown>;
	text(): Promise<string>;
};

export type FetchInit = {
	method?: string;
	headers?: Record<string, string>;
	body?: string;
	signal?: AbortSignal | null;
};

export type FetchFn = (
	input: string,
	init?: FetchInit
) => Promise<FetchResponse>;

export interface RPCClientOptions {
	prefix?: string;
	headers?: Record<string, string>;
	bearerToken?: string;
	timeoutMs?: number;
	fetchFn?: FetchFn;
}

export class RPCClient {
	private readonly baseURL: string;
	private readonly prefix: string;
	private readonly headers: Record<string, string>;
	private readonly bearerToken: string;
	private readonly timeoutMs?: number;
	private readonly fetchFn: FetchFn;

	constructor(baseURL: string, options: RPCClientOptions = {}) {
		this.baseURL = RPCClient.normalizeBaseURL(baseURL);
		this.prefix = RPCClient.normalizePrefix(options.prefix ?? "/rpc");
		this.headers = { ...(options.headers ?? {}) };
		this.bearerToken = options.bearerToken ?? "";
		this.timeoutMs = options.timeoutMs;
		this.fetchFn =
			options.fetchFn ??
			(async (input, init) =>
				(fetch(input, init as unknown as RequestInit) as unknown as FetchResponse));
	}

	private static normalizeBaseURL(baseURL: string): string {
		const trimmed = baseURL.trim();
		const withScheme = trimmed.includes("://") ? trimmed : `http://${trimmed}`;
		return withScheme.replace(/\/+$/g, "");
	}

	private static normalizePrefix(prefix: string): string {
		const trimmed = prefix.trim();
		if (trimmed === "") {
			return "";
		}
		const withSlash = trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
		return withSlash.replace(/\/+$/g, "");
	}

	private buildURL(path: string): string {
		if (this.prefix) {
			return `${this.baseURL}${this.prefix}/${path}`;
		}
		return `${this.baseURL}/${path}`;
	}

	private async request(path: string, payload?: unknown): Promise<unknown> {
		const headers: Record<string, string> = {
			"Content-Type": "application/json",
			...this.headers,
		};
		if (this.bearerToken && !hasHeader(this.headers, "Authorization")) {
			headers.Authorization = `Bearer ${this.bearerToken}`;
		}

		const controller = this.timeoutMs ? new AbortController() : undefined;
		const timeout = this.timeoutMs
			? setTimeout(() => controller?.abort(), this.timeoutMs)
			: undefined;
		try {
			const response = await this.fetchFn(this.buildURL(path), {
				method: "POST",
				headers,
				body: payload === undefined ? undefined : JSON.stringify(payload),
				signal: controller?.signal,
			});

			if (!response.ok) {
				let parsed: RPCError | undefined;
				try {
					parsed = (await response.json()) as RPCError;
				} catch {
					parsed = undefined;
				}
				if (parsed && parsed.type) {
					this.raiseError(parsed);
				}
				throw new RPCErrorException({
					type: "custom",
					message: `rpc error: status ${response.status}`,
				});
			}

			const body = await response.text();
			if (body.trim() === "") {
				return undefined;
			}
			return JSON.parse(body);
		} finally {
			if (timeout) {
				clearTimeout(timeout);
			}
		}
	}

	private raiseError(error: RPCError): never {
		const excType = ERROR_EXCEPTIONS[error.type];
		if (excType) {
			throw new excType(error);
		}
		throw new RPCErrorException(error);
	}
	async testEmpty(): Promise<EmptyModel> {
		const payload = undefined;
		const res = (await this.request("test_empty", payload)) as TestEmptyResult;
		return res.empty;
	}
	async testNoReturn(): Promise<void> {
		const payload = undefined;
		await this.request("test_no_return", payload);
	}
	async testBasic(params: TestBasicParams): Promise<TextModel> {
		const payload = TestBasicParamsSchema.parse(params);
		const res = (await this.request("test_basic", payload)) as TestBasicResult;
		return res.text;
	}
	async testListMap(params: TestListMapParams): Promise<NestedModel> {
		const payload = TestListMapParamsSchema.parse(params);
		const res = (await this.request("test_list_map", payload)) as TestListMapResult;
		return res.nested;
	}
	async testOptional(params: TestOptionalParams): Promise<FlagsModel> {
		const payload = TestOptionalParamsSchema.parse(params);
		const res = (await this.request("test_optional", payload)) as TestOptionalResult;
		return res.flags;
	}
	async testValidationError(params: TestValidationErrorParams): Promise<TextModel> {
		const payload = TestValidationErrorParamsSchema.parse(params);
		const res = (await this.request("test_validation_error", payload)) as TestValidationErrorResult;
		return res.text;
	}
	async testUnauthorizedError(): Promise<EmptyModel> {
		const payload = undefined;
		const res = (await this.request("test_unauthorized_error", payload)) as TestUnauthorizedErrorResult;
		return res.empty;
	}
	async testForbiddenError(): Promise<EmptyModel> {
		const payload = undefined;
		const res = (await this.request("test_forbidden_error", payload)) as TestForbiddenErrorResult;
		return res.empty;
	}
	async testNotImplementedError(): Promise<EmptyModel> {
		const payload = undefined;
		const res = (await this.request("test_not_implemented_error", payload)) as TestNotImplementedErrorResult;
		return res.empty;
	}
	async testCustomError(): Promise<EmptyModel> {
		const payload = undefined;
		const res = (await this.request("test_custom_error", payload)) as TestCustomErrorResult;
		return res.empty;
	}
	async testMapReturn(): Promise<Record<string, TextModel>> {
		const payload = undefined;
		const res = (await this.request("test_map_return", payload)) as TestMapReturnResult;
		return res.result;
	}
	async testJson(params: TestJsonParams): Promise<any> {
		const payload = TestJsonParamsSchema.parse(params);
		const res = (await this.request("test_json", payload)) as TestJsonResult;
		return res.json;
	}
	async testRaw(params: TestRawParams): Promise<any> {
		const payload = TestRawParamsSchema.parse(params);
		const res = (await this.request("test_raw", payload)) as TestRawResult;
		return res.raw;
	}
	async testMixedPayload(params: TestMixedPayloadParams): Promise<PayloadModel> {
		const payload = TestMixedPayloadParamsSchema.parse(params);
		const res = (await this.request("test_mixed_payload", payload)) as TestMixedPayloadResult;
		return res.payload;
	}
}

function hasHeader(headers: Record<string, string>, name: string): boolean {
	const target = name.toLowerCase();
	for (const key of Object.keys(headers)) {
		if (key.toLowerCase() === target) {
			return true;
		}
	}
	return false;
}
