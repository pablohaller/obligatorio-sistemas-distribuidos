import { FormField, defaultStringFormField } from "@/app/utils/helpers";

export interface PasswordChecks {
  hasLength: boolean | undefined;
  hasNumbers: boolean | undefined;
  hasSymbols: boolean | undefined;
  hasUppercase: boolean | undefined;
}

export interface FormData {
  name: FormField<string>;
  email: FormField<string>;
  password: FormField<string>;
  passwordCheck: FormField<string>;
}

export const defaultPasswordChecks = {
  hasLength: undefined,
  hasNumbers: undefined,
  hasSymbols: undefined,
  hasUppercase: undefined,
};

export const defaultFormData: FormData = {
  name: { ...defaultStringFormField, required: true, label: "Nombre completo" },
  email: { ...defaultStringFormField, required: true, label: "E-mail" },
  password: { ...defaultStringFormField, required: true, label: "Contraseña" },
  passwordCheck: {
    ...defaultStringFormField,
    required: true,
    label: "Repetir contraseña",
  },
};

export const getPasswordCheks = (value: string) => {
  const specialChars = /[`!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/;
  const numbers = /\d/;
  const uppercase = /[A-Z]/;

  return {
    hasLength: value.length > 8,
    hasNumbers: numbers.test(value),
    hasSymbols: specialChars.test(value),
    hasUppercase: uppercase.test(value),
  };
};
