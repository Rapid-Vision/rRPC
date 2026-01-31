// THIS CODE IS GENERATED


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
