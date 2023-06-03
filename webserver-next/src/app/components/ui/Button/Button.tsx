import React from "react";
import { twMerge } from "tailwind-merge";

interface Props {
  children: React.ReactNode;
  variant?: "default" | "contained";
  fullWidth?: boolean;
  noPadding?: boolean;
  onClick: (e?: unknown) => unknown;
}

const VARIANTS = {
  default: "font-light text-sky-500 hover:underline",
  contained:
    "font-bold text-white bg-sky-500 hover:bg-sky-600 active:bg-sky-700",
};

const Button = ({
  children,
  variant = "default",
  fullWidth,
  noPadding,
  onClick,
}: Props) => {
  return (
    <button
      className={twMerge(
        "p-2 rounded-md font-rubik",
        VARIANTS[variant],
        fullWidth && "w-full",
        noPadding && "p-0"
      )}
      onClick={onClick}
    >
      {children}
    </button>
  );
};

export default Button;
