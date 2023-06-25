import { ChangeEvent } from "react";
import { twMerge } from "tailwind-merge";

interface Props {
  label: string;
  type?: string;
  required?: boolean;
  placeholder?: string;
  disabled?: boolean;
  error?: string;
  value: string;
  onChange: (event: ChangeEvent<HTMLInputElement>) => void;
  onBlur?: () => void;
  onClick?: () => void;
}

const Input = ({
  label,
  required,
  placeholder,
  error,
  disabled,
  value,
  type = "text",
  onChange,
  onBlur,
}: Props) => {
  return (
    <div
      className={twMerge(
        "w-full mb-1",
        disabled && "opacity-30 pointer-events-none"
      )}
    >
      {label && (
        <div className="font-rubik font-light text-base">
          {label} {required && <span className="text-red-500">*</span>}
        </div>
      )}
      <input
        type={type}
        className={twMerge(
          `w-full border rounded-md p-2 bg-gray-100 focus:bg-white focus:outline-none focus-visible:ring-2 focus-visible:ring-sky-500 text-sm placeholder:text-md text-md`,
          error && "border-red-500"
        )}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        maxLength={50}
      />
      <div className="text-xs text-red-500">&nbsp;{error}</div>
    </div>
  );
};

export default Input;
