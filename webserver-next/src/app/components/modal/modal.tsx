import Button from "../ui/Button/Button";

interface Props {
  title?: string;
  children?: React.ReactNode;
  onConfirm?: () => void;
  confirmButtonText?: string;
  showCancelButton?: boolean;
  hideConfirmButton?: boolean;
  cancelButtonText?: string;
  onCancel?: () => void;
}

const Modal = ({
  title,
  children,
  onConfirm,
  confirmButtonText = "Confirmar",
  showCancelButton,
  cancelButtonText = "Cancelar",
  hideConfirmButton,
  onCancel,
}: Props) => {
  return (
    <div className="z-[1000] fixed top-0 left-0 w-screen h-screen bg-white/50 backdrop-blur-sm grid place-items-center ">
      <div className="bg-white p-4 rounded-xl drop-shadow-md md:min-w-1/5 flex flex-col justify-between">
        <div>
          <div className="text-xl font-semibold mb-2">{title}</div>
          <div>{children}</div>
        </div>
        <div className="flex flex-col mt-4">
          {!hideConfirmButton && (
            <Button onClick={onConfirm} variant="danger">
              {confirmButtonText}
            </Button>
          )}
          {showCancelButton && (
            <Button onClick={onCancel}>{cancelButtonText}</Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default Modal;
