import { useSearchParams } from "react-router-dom";

const FailPaymentPage = () => {
  const [searchParams] = useSearchParams();

  return (
    <div>
      <h1>결제 실패</h1>
      <div>{`사유: ${searchParams.get("message")}`}</div>
    </div>
  );
};

export default FailPaymentPage;
