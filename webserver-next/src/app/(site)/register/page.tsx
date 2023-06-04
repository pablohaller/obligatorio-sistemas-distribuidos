"use client";
import Button from "@/app/components/ui/Button/Button";
import Input from "@/app/components/ui/Input/Input";
import { LOGIN_PATH } from "@/app/constants/routes";
import {
  IconCircleCheck,
  IconLayoutDashboard,
  IconLoader3,
} from "@tabler/icons-react";
import { useRouter } from "next/navigation";
import { useCallback, useState } from "react";
import { toast } from "react-toastify";
import {
  FormData,
  PasswordChecks,
  defaultFormData,
  defaultPasswordChecks,
  getPasswordCheks,
} from "./helpers";
import PasswordChecksList from "./PasswordChecksList";
import {
  checkPayloadEmptyFiels,
  getFormDataErrors,
  getFormDataPayload,
  handleFormData,
} from "@/app/utils/helpers";
import { twMerge } from "tailwind-merge";

export default function Page() {
  const router = useRouter();
  const [formData, setFormData] = useState<FormData>(defaultFormData);
  const [passwordChecks, setPasswordChecks] = useState<PasswordChecks>(
    defaultPasswordChecks
  );
  const [loading, setLoading] = useState<boolean>(false);
  const [successMessage, setSuccessMessage] = useState<boolean>(false);

  const handleFrom = useCallback(
    (key: string, attribute: string, value: unknown) => {
      setFormData(handleFormData(formData, key, attribute, value));
    },
    [formData]
  );

  const handleRegister = useCallback(async () => {
    setLoading(true);
    const payload = getFormDataPayload(formData);
    const emptyFields = checkPayloadEmptyFiels(payload);

    if (emptyFields.length) {
      setFormData(getFormDataErrors(formData, emptyFields));
      setLoading(false);
      return;
    }

    if (payload?.password !== payload?.passwordCheck) {
      toast.error("Las contraseñas no coinciden", { theme: "colored" });
      setLoading(false);
      return;
    }

    const request = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/users`,
      {
        method: "POST",
        body: JSON.stringify(payload),
      }
    );

    const data = await request.json();

    if (data?.error) {
      toast.error(data?.error, { theme: "colored" });
      setLoading(false);
      return;
    }

    setLoading(false);
    setSuccessMessage(true);
  }, [formData]);

  const {
    name: {
      label: nameLabel,
      value: nameValue,
      error: nameError,
      required: nameRequired,
    },
    email: {
      label: emailLabel,
      value: emailValue,
      error: emailError,
      required: emailRequired,
    },
    password: {
      label: passwordLabel,
      value: passwordValue,
      error: passwordError,
      required: passwordRequired,
    },
    passwordCheck: {
      label: passwordCheckLabel,
      value: passwordCheckValue,
      error: passwordCheckError,
      required: passwordCheckRequired,
    },
  } = formData;

  return (
    <div className="h-full grid place-items-center p-2">
      <div className="  bg-white p-4 rounded-xl drop-shadow-md md:w-2/5 relative">
        {loading && (
          <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-10 flex flex-col items-center">
            <IconLoader3 className="h-20 w-20 text-sky-500 animate-spin" />
            <div className="text-3xl">Cargando</div>
          </div>
        )}
        {successMessage && (
          <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-10 flex flex-col items-center text-center">
            <IconCircleCheck className="h-20 w-20 text-sky-500" />
            <div className="text-3xl">Usuario registrado con éxito</div>
            <Button
              variant="default"
              fullWidth
              noPadding
              onClick={() => router.push(LOGIN_PATH)}
            >
              Haga click aquí para iniciar sesión
            </Button>
          </div>
        )}
        <div
          className={twMerge(
            (loading || successMessage) && "blur-sm pointer-events-none"
          )}
        >
          <IconLayoutDashboard className="h-8 w-8 md:h-12 md:w-12 text-sky-500 m-auto" />
          <div className="font-rubik font-500 text-sky-500 text-center p-5 text-xl md:text-3xl">
            NEXTJS APP
          </div>
          <Input
            label={nameLabel}
            placeholder={nameLabel}
            required={nameRequired}
            value={nameValue}
            error={nameError}
            onChange={(e) => handleFrom("name", "value", e.target.value)}
          />
          <Input
            label={emailLabel}
            placeholder={emailLabel}
            required={emailRequired}
            value={emailValue}
            error={emailError}
            onChange={(e) => handleFrom("email", "value", e.target.value)}
          />
          <Input
            label={passwordLabel}
            placeholder={passwordLabel}
            type="password"
            required={passwordRequired}
            value={passwordValue}
            error={passwordError}
            onChange={(e) => {
              handleFrom("password", "value", e.target.value);
              setPasswordChecks(getPasswordCheks(e.target.value));
            }}
          />
          <Input
            label={passwordCheckLabel}
            placeholder={passwordCheckLabel}
            type="password"
            required={passwordCheckRequired}
            value={passwordCheckValue}
            error={passwordCheckError}
            onChange={(e) =>
              handleFrom("passwordCheck", "value", e.target.value)
            }
          />
          <div className="text-sm mb-4">
            La contraseña deben tener al menos:
            <PasswordChecksList passwordChecks={passwordChecks} />
          </div>
          <Button variant="contained" fullWidth onClick={handleRegister}>
            Registrarse
          </Button>
          <div className="mt-4">
            <Button
              variant="default"
              fullWidth
              noPadding
              onClick={() => router.push(LOGIN_PATH)}
            >
              Iniciar sesión
            </Button>
            <Button
              variant="default"
              noPadding
              fullWidth
              onClick={() => {
                toast.warn("No implementado aún", { theme: "colored" });
              }}
            >
              Recuperar contraseña
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
