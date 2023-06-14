export interface FormField<T> {
  label: string;
  value: T;
  error?: string;
  required?: boolean;
}

export const defaultStringFormField = {
  label: "",
  value: "",
  error: "",
};

export const getFormDataPayload = (formData: Object) => {
  const payload: any = {};
  Object.entries(formData).forEach(([key, value]) => {
    payload[key] = value.value;
  });
  return payload;
};

export const checkPayloadEmptyFiels = (payload: Object) =>
  Object.entries(payload)
    .filter((entry) => !entry[1])
    .map(([key]) => key);

export const handleFormData = (
  formData: any,
  key: string,
  attribute: string,
  value: unknown
) => ({
  ...formData,
  [key]: {
    ...formData[key as keyof any],
    [attribute]: value,
    ...(attribute !== "error" ? { error: "" } : {}),
  },
});

export const getFormDataErrors = (formData: Object, fields: string[]) => {
  let formDataErrors: any = { ...formData };
  fields.forEach((field) => {
    const label = `${
      formDataErrors[field as keyof any]?.label
    } no puede ser vacío`;
    formDataErrors = handleFormData(formDataErrors, field, "error", label);
  });
  return formDataErrors;
};
