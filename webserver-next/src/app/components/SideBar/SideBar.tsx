"use client";
import {
  IconEngine,
  IconHome,
  IconLogout,
  IconMap,
  IconMenu2,
  IconReport,
  IconUserCircle,
  IconX,
} from "@tabler/icons-react";
import { Session } from "next-auth";
import { signOut } from "next-auth/react";
import { useCallback, useEffect, useState } from "react";
import Modal from "../modal/modal";
import Link from "next/link";
import { twMerge } from "tailwind-merge";
import { usePathname } from "next/navigation";
import LoadingSpinner from "../LoadingSpinner/LoadingSpinner";
import {
  DASHBOARD_MAP_PATH,
  DASHBOARD_MEASUREMENTS_PATH,
  DASHBOARD_PATH,
  DASHBOARD_REPORTS_PATH,
} from "@/app/constants/routes";

interface Props {
  session: Session | null;
}

interface ActivePathname {
  isHomeActive: boolean;
  isReportActive: boolean;
  isMapActive: boolean;
  isMeasurementsActive: boolean;
}

const getActivePathName = (pathname: string | undefined): ActivePathname => ({
  isHomeActive: pathname === DASHBOARD_PATH || false,
  isReportActive: pathname?.includes(DASHBOARD_REPORTS_PATH) || false,
  isMapActive: pathname?.includes(DASHBOARD_MAP_PATH) || false,
  isMeasurementsActive:
    pathname?.includes(DASHBOARD_MEASUREMENTS_PATH) || false,
});

const SideBar = ({ session }: Props) => {
  const clientPathName = usePathname();
  const [loadingSignOut, setLoadingSignOut] = useState(false);
  const [activePathname, setActivePathname] = useState<ActivePathname>(
    getActivePathName(clientPathName)
  );
  const [showLogOut, setShowLogOut] = useState<boolean>(false);
  const handleShowLogOut = useCallback(
    () => setShowLogOut(!showLogOut),
    [showLogOut]
  );

  const [showMobileMenu, setShowMobileMenu] = useState<boolean>(false);
  const handleMobileMenu = useCallback(
    () => setShowMobileMenu(!showMobileMenu),
    [showMobileMenu]
  );
  const handleSignOut = useCallback(() => {
    signOut();
    setLoadingSignOut(true);
  }, []);

  useEffect(() => {
    setActivePathname(getActivePathName(clientPathName));
  }, [clientPathName]);

  const { isReportActive, isMapActive, isMeasurementsActive, isHomeActive } =
    activePathname;

  return (
    <>
      <div className="md:w-1/5 bg-white/60 md:shadow-2xl backdrop-blur-sm relative z-20 md:static">
        <div className="flex justify-between items-center py-4 px-4 border-b border-gray-300 shadow-sm min-w-full flex-row-reverse md:flex-row">
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
            className=" text-sky-500 cursor-pointer flex-shrink-0 hidden md:block"
            onClick={handleShowLogOut}
          />
          {
            <div
              className=" text-sky-500 cursor-pointer flex-shrink-0 block md:hidden"
              onClick={handleMobileMenu}
            >
              {showMobileMenu ? <IconX /> : <IconMenu2 />}
            </div>
          }
        </div>
        <div
          className={twMerge(
            "z-10 p-2 font-rubik h-screen md:h-0 absolute w-full md:bg-transparent  bg-white/90 shadow-2xl backdrop-blur-md md:bg-none md:backdrop-blur-none md:shadow-none",
            !showMobileMenu && "hidden md:block",
            showMobileMenu && "h-screen md:h-0"
          )}
        >
          <Link
            className={twMerge(
              "p-4 flex items-center hover:text-sky-400",
              isHomeActive &&
                "text-sky-500  bg-sky-100 rounded-xl hover:text-sky-500"
            )}
            href={DASHBOARD_PATH}
          >
            <IconHome className="mr-2" />
            <span>Inicio</span>
          </Link>
          <Link
            className={twMerge(
              "p-4 flex items-center hover:text-sky-400",
              isMapActive &&
                "text-sky-500  bg-sky-100 rounded-xl hover:text-sky-500"
            )}
            href={DASHBOARD_MAP_PATH}
          >
            <IconMap className="mr-2" />
            <span>Mapa</span>
          </Link>
          <Link
            className={twMerge(
              "p-4 flex items-center hover:text-sky-400",
              isReportActive &&
                "text-sky-500 bg-sky-100 rounded-xl  hover:text-sky-500"
            )}
            href={DASHBOARD_REPORTS_PATH}
          >
            <IconReport className="mr-2" />
            <span>Reportes</span>
          </Link>
          <Link
            className={twMerge(
              "p-4 flex items-center hover:text-sky-400",
              isMeasurementsActive &&
                "text-sky-500 bg-sky-100 rounded-xl  hover:text-sky-500"
            )}
            href={DASHBOARD_MEASUREMENTS_PATH}
          >
            <IconEngine className="mr-2 flex-shrink-0" />
            <span>Mediciones</span>
          </Link>
          <div
            className="p-4 flex items-center cursor-pointer hover:text-sky-400"
            onClick={handleShowLogOut}
          >
            <IconLogout className=" flex-shrink-0 mr-2" />
            Cerrar sesión
          </div>
        </div>
      </div>
      {loadingSignOut && (
        <Modal hideConfirmButton>
          <div className="h-32 w-36">
            {loadingSignOut && <LoadingSpinner />}
          </div>
        </Modal>
      )}
      {showLogOut && !loadingSignOut && (
        <Modal
          title="Cerrar sesión"
          onConfirm={handleSignOut}
          showCancelButton
          onCancel={handleShowLogOut}
        >
          <span>¿Seguro desea cerrar sesión?</span>
        </Modal>
      )}
    </>
  );
};

export default SideBar;
