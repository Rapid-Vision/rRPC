// THIS CODE IS GENERATED

export type RPCErrorType =
	| "custom"
	| "validation"
	| "input"
	| "unauthorized"
	| "forbidden"
	| "not_implemented";

export interface RPCError {
	type: RPCErrorType;
	message: string;
}

export class RPCErrorException extends Error {
	readonly error: RPCError;

	constructor(error: RPCError) {
		super(error.message);
		this.error = error;
	}
}

export class CustomRPCError extends RPCErrorException {}
export class ValidationRPCError extends RPCErrorException {}
export class InputRPCError extends RPCErrorException {}
export class UnauthorizedRPCError extends RPCErrorException {}
export class ForbiddenRPCError extends RPCErrorException {}
export class NotImplementedRPCError extends RPCErrorException {}

const ERROR_EXCEPTIONS: Record<string, typeof RPCErrorException> = {
	custom: CustomRPCError,
	validation: ValidationRPCError,
	input: InputRPCError,
	unauthorized: UnauthorizedRPCError,
	forbidden: ForbiddenRPCError,
	not_implemented: NotImplementedRPCError,
};
export interface EmptyModel {
}
export interface TextModel {
	title?: string | null;
	body: string;
}
export interface FlagsModel {
	enabled: boolean;
	retries: number;
	labels: Array<string>;
	meta: Record<string, string>;
}
export interface NestedModel {
	text: TextModel;
	flags?: FlagsModel | null;
	items: Array<TextModel>;
	lookup: Record<string, TextModel>;
}
export interface PayloadModel {
	data: any;
	raw_data: any;
}
export interface TestEmptyResult {
	empty: EmptyModel;
}
export interface TestBasicParams {
	text: TextModel;
	flag: boolean;
	count: number;
	note?: string | null;
}
export interface TestBasicResult {
	text: TextModel;
}
export interface TestListMapParams {
	texts: Array<TextModel>;
	flags: Record<string, string>;
}
export interface TestListMapResult {
	nested: NestedModel;
}
export interface TestOptionalParams {
	text?: TextModel | null;
	flag?: boolean | null;
}
export interface TestOptionalResult {
	flags: FlagsModel;
}
export interface TestValidationErrorParams {
	text: TextModel;
}
export interface TestValidationErrorResult {
	text: TextModel;
}
export interface TestUnauthorizedErrorResult {
	empty: EmptyModel;
}
export interface TestForbiddenErrorResult {
	empty: EmptyModel;
}
export interface TestNotImplementedErrorResult {
	empty: EmptyModel;
}
export interface TestCustomErrorResult {
	empty: EmptyModel;
}
export interface TestMapReturnResult {
	result: Record<string, TextModel>;
}
export interface TestJsonParams {
	data: any;
}
export interface TestJsonResult {
	json: any;
}
export interface TestRawParams {
	payload: any;
}
export interface TestRawResult {
	raw: any;
}
export interface TestMixedPayloadParams {
	payload: PayloadModel;
}
export interface TestMixedPayloadResult {
	payload: PayloadModel;
}

export interface RPCClientOptions {
	prefix?: string;
	headers?: Record<string, string>;
	bearerToken?: string;
	timeoutMs?: number;
	fetchFn?: typeof fetch;
}

export class RPCClient {
	private readonly baseURL: string;
	private readonly prefix: string;
	private readonly headers: Record<string, string>;
	private readonly bearerToken: string;
	private readonly timeoutMs?: number;
	private readonly fetchFn: typeof fetch;

	constructor(baseURL: string, options: RPCClientOptions = {}) {
		this.baseURL = RPCClient.normalizeBaseURL(baseURL);
		this.prefix = RPCClient.normalizePrefix(options.prefix ?? "/rpc");
		this.headers = { ...(options.headers ?? {}) };
		this.bearerToken = options.bearerToken ?? "";
		this.timeoutMs = options.timeoutMs;
		this.fetchFn = options.fetchFn ?? fetch;
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
		const payload = params;
		const res = (await this.request("test_basic", payload)) as TestBasicResult;
		return res.text;
	}
	async testListMap(params: TestListMapParams): Promise<NestedModel> {
		const payload = params;
		const res = (await this.request("test_list_map", payload)) as TestListMapResult;
		return res.nested;
	}
	async testOptional(params: TestOptionalParams): Promise<FlagsModel> {
		const payload = params;
		const res = (await this.request("test_optional", payload)) as TestOptionalResult;
		return res.flags;
	}
	async testValidationError(params: TestValidationErrorParams): Promise<TextModel> {
		const payload = params;
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
		const payload = params;
		const res = (await this.request("test_json", payload)) as TestJsonResult;
		return res.json;
	}
	async testRaw(params: TestRawParams): Promise<any> {
		const payload = params;
		const res = (await this.request("test_raw", payload)) as TestRawResult;
		return res.raw;
	}
	async testMixedPayload(params: TestMixedPayloadParams): Promise<PayloadModel> {
		const payload = params;
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
