import {
  IconCircleCheck,
  IconCircleDotted,
  IconCircleX,
} from "@tabler/icons-react";
import { PasswordChecks } from "./helpers";

const getIcon = (checkState: boolean | undefined) => {
  if (checkState === true) {
    return <IconCircleCheck className="mr-2 text-sky-500" />;
  }
  if (checkState === false) {
    return <IconCircleX className="mr-2 text-red-500" />;
  }
  return <IconCircleDotted className="mr-2" />;
};

interface Props {
  passwordChecks: PasswordChecks;
}
const PasswordChecksList = ({ passwordChecks }: Props) => {
  const { hasLength, hasNumbers, hasSymbols, hasUppercase } = passwordChecks;

  return (
    <div className="grid grid-cols-2">
      <div className="flex items-center">
        {getIcon(hasLength)}
        <span>8 caracteres</span>
      </div>
      <div className="flex items-center">
        {getIcon(hasNumbers)}
        <span>1 n&uacute;mero</span>
      </div>
      <div className="flex items-center">
        {getIcon(hasSymbols)}
        <span>1 s&iacute;mbolo</span>
      </div>
      <div className="flex items-center">
        {getIcon(hasUppercase)}
        <span>1 may&uacute;scula</span>
      </div>
    </div>
  );
};

export default PasswordChecksList;
