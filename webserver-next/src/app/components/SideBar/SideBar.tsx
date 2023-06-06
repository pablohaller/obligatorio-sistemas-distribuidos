"use client";
import {
  IconLogout,
  IconMap,
  IconReport,
  IconUserCircle,
} from "@tabler/icons-react";
import { Session } from "next-auth";
import { signOut } from "next-auth/react";
import React, { useCallback, useEffect, useState } from "react";
import Modal from "../modal/modal";
import Link from "next/link";
import { twMerge } from "tailwind-merge";
import { usePathname } from "next/navigation";

interface Props {
  session: Session | null;
}

interface ActivePathname {
  isReportActive: boolean;
  isMapActive: boolean;
}

const getActivePathName = (pathname: string | undefined): ActivePathname => ({
  isReportActive: pathname?.includes("/dashboard/reports") || false,
  isMapActive: pathname?.includes("/dashboard/map") || false,
});

const SideBar = ({ session }: Props) => {
  const clientPathName = usePathname();
  const [showLogOut, setShowLogOut] = useState<boolean>(false);
  const [activePathname, setActivePathname] = useState<ActivePathname>(
    getActivePathName(clientPathName)
  );
  const handleShowLogOut = useCallback(
    () => setShowLogOut(!showLogOut),
    [showLogOut]
  );

  useEffect(() => {
    setActivePathname(getActivePathName(clientPathName));
  }, [clientPathName]);

  const { isReportActive, isMapActive } = activePathname;

  return (
    <>
      <div className="w-1/5 bg-white/60 shadow-2xl backdrop-blur-sm">
        <div className="flex justify-between items-center py-4 px-4 border-b border-gray-300 shadow-sm min-w-full">
          <div className="flex items-center min-w-0 ">
            <IconUserCircle className="h-10 w-10 mr-2 text-sky-500 flex-shrink-0" />
            <div className="flex flex-col min-w-0">
              <div className="overflow-hidden text-ellipsis whitespace-nowrap text-xl font-rubik ">
                {session?.user?.name}
              </div>
              <div className="overflow-hidden text-ellipsis whitespace-nowrap text-xs font-rubik ">
                {session?.user?.email}
              </div>
            </div>
          </div>
          <IconLogout
            className=" text-sky-500 cursor-pointer flex-shrink-0"
            // onClick={() => signOut()}
            onClick={handleShowLogOut}
          />
        </div>
        <div className="p-2">
          <Link
            className={twMerge(
              "p-4 text-xl flex items-center",
              isReportActive && "text-sky-500"
            )}
            href="/dashboard/reports"
          >
            <IconReport className="mr-2" />
            <span>Reportes</span>
          </Link>
          <Link
            className={twMerge(
              "p-4 text-xl flex items-center",
              isMapActive && "text-sky-500"
            )}
            href="/dashboard/map"
          >
            <IconMap className="mr-2" />
            <span>Mapa</span>
          </Link>
        </div>
      </div>
      {showLogOut && (
        <Modal
          title="Cerrar sesión"
          onConfirm={signOut}
          showCancelButton
          onCancel={handleShowLogOut}
        >
          ¿Seguro desea cerrar sesión?
        </Modal>
      )}
    </>
  );
};

export default SideBar;
