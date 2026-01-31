// THIS CODE IS GENERATED


import { z } from "zod";
export interface EmptyModel {
}

export const EmptyModelSchema = z.object({
});
export interface TextModel {
	title?: string | null;
	body: string;
}

export const TextModelSchema = z.object({
	title: z.union([z.string(), z.null()]).optional(),
	body: z.string(),
});
export interface FlagsModel {
	enabled: boolean;
	retries: number;
	labels: Array<string>;
	meta: Record<string, string>;
}

export const FlagsModelSchema = z.object({
	enabled: z.boolean(),
	retries: z.number().int(),
	labels: z.array(z.string()),
	meta: z.record(z.string(), z.string()),
});
export interface NestedModel {
	text: TextModel;
	flags?: FlagsModel | null;
	items: Array<TextModel>;
	lookup: Record<string, TextModel>;
}

export const NestedModelSchema = z.object({
	text: z.lazy(() => TextModelSchema),
	flags: z.union([z.lazy(() => FlagsModelSchema), z.null()]).optional(),
	items: z.array(z.lazy(() => TextModelSchema)),
	lookup: z.record(z.string(), z.lazy(() => TextModelSchema)),
});
export interface PayloadModel {
	data: any;
	raw_data: any;
}

export const PayloadModelSchema = z.object({
	data: z.any(),
	raw_data: z.any(),
});
export interface TestEmptyResult {
	empty: EmptyModel;
}
export interface TestBasicParams {
	text: TextModel;
	flag: boolean;
	count: number;
	note?: string | null;
}

export const TestBasicParamsSchema = z.object({
	text: z.lazy(() => TextModelSchema),
	flag: z.boolean(),
	count: z.number().int(),
	note: z.union([z.string(), z.null()]).optional(),
});
export interface TestBasicResult {
	text: TextModel;
}
export interface TestListMapParams {
	texts: Array<TextModel>;
	flags: Record<string, string>;
}

export const TestListMapParamsSchema = z.object({
	texts: z.array(z.lazy(() => TextModelSchema)),
	flags: z.record(z.string(), z.string()),
});
export interface TestListMapResult {
	nested: NestedModel;
}
export interface TestOptionalParams {
	text?: TextModel | null;
	flag?: boolean | null;
}

export const TestOptionalParamsSchema = z.object({
	text: z.union([z.lazy(() => TextModelSchema), z.null()]).optional(),
	flag: z.union([z.boolean(), z.null()]).optional(),
});
export interface TestOptionalResult {
	flags: FlagsModel;
}
export interface TestValidationErrorParams {
	text: TextModel;
}

export const TestValidationErrorParamsSchema = z.object({
	text: z.lazy(() => TextModelSchema),
});
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

export const TestJsonParamsSchema = z.object({
	data: z.any(),
});
export interface TestJsonResult {
	json: any;
}
export interface TestRawParams {
	payload: any;
}

export const TestRawParamsSchema = z.object({
	payload: z.any(),
});
export interface TestRawResult {
	raw: any;
}
export interface TestMixedPayloadParams {
	payload: PayloadModel;
}

export const TestMixedPayloadParamsSchema = z.object({
	payload: z.lazy(() => PayloadModelSchema),
});
export interface TestMixedPayloadResult {
	payload: PayloadModel;
}
