import { IconEye } from "@tabler/icons-react";
import React, { useRef } from "react";
import { twMerge } from "tailwind-merge";

interface Props {
  sector: string;
  isSelected: boolean;
  onClick: () => void;
}

const SectorListItem = ({ sector, isSelected, onClick }: Props) => {
  const divRef = useRef<HTMLDivElement>(null);

  return (
    <div
      ref={divRef}
      className={twMerge(
        "border rounded-xl p-2 flex  flex-col bg-white mb- border-l-8",
        isSelected ? "border-l-sky-500" : "border-l-gray-300"
      )}
    >
      <div>{sector}</div>
      <div
        className="text-sm flex items-center   hover:text-sky-500 cursor-pointer"
        onClick={() => {
          onClick();
          divRef?.current?.scrollIntoView({
            behavior: "smooth",
          });
        }}
      >
        {isSelected ? (
          <span>&nbsp;</span>
        ) : (
          <>
            <IconEye className="mr-2" />
            <span className="hover:underline">Ver sector</span>
          </>
        )}
      </div>
    </div>
  );
};

export default SectorListItem;
