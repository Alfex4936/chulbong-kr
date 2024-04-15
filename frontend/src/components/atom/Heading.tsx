type Props = { title: string; subTitle?: string };

const Heading = ({ title, subTitle }: Props) => {
  return (
    <div className="font-medium text-2xl text-center p-10 mo:text-lg">
      <div>{title}</div>
      {subTitle && <div className="text-sm text-grey-dark mo:text-xs">({subTitle})</div>}
    </div>
  );
};

export default Heading;
