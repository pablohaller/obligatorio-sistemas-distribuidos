import React, { ChangeEvent, ChangeEventHandler, SetStateAction } from "react";
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
}: Props) => {
  return (
    <div
      className={twMerge(
        "w-full",
        disabled && "opacity-30 pointer-events-none"
      )}
    >
      {label && (
        <div>
          {label} {required && <span className="text-red-500">*</span>}
        </div>
      )}
      <input
        type={type}
        className="w-full border rounded-md p-2 bg-gray-100 focus:bg-white focus:outline-none focus-visible:ring-2 focus-visible:ring-sky-500"
        placeholder={placeholder}
        value={value}
        onChange={onChange}
      />
      <div className="text-xs text-red-500">&nbsp;{error}</div>
    </div>
  );
};

export default Input;
