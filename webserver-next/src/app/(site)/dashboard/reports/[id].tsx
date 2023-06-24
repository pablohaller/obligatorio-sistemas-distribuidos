interface Props {
  params: { id: string };
}

const Page = ({ params }: Props) => {
  console.log("params", params);
  return <div>My Id: {params?.id}</div>;
};

export default Page;
