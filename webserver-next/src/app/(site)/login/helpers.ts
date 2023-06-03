import { FormField, defaultStringFormField } from "@/app/utils/helpers";

export interface FormData {
  email: FormField<string>;
  password: FormField<string>;
}

export const defaultFormData: FormData = {
  email: { ...defaultStringFormField, required: true, label: "E-mail" },
  password: { ...defaultStringFormField, required: true, label: "Contraseña" },
};
