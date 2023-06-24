interface Props {
  children: React.ReactNode;
}

const Layout = async ({ children }: Props) => {
  return <>{children}</>;
};

export default Layout;
