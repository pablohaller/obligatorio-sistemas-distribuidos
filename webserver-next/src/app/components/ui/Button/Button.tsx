import { twMerge } from "tailwind-merge";

interface Props {
  children: React.ReactNode;
  variant?: "default" | "contained" | "danger";
  fullWidth?: boolean;
  noPadding?: boolean;
  onClick: (() => void) | undefined;
}

const VARIANTS = {
  default: "font-light text-sky-500 hover:underline",
  contained:
    "font-bold text-white bg-sky-500 hover:bg-sky-600 active:bg-sky-700",
  danger: "font-bold text-white bg-red-500 hover:bg-red-600 active:bg-red-700",
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
