import { useSearchParams } from "react-router-dom";

const SuccessPaymentPage = () => {
  const [searchParams] = useSearchParams();

  return (
    <div>
      <h1>결제 성공</h1>
      <div>{`주문 아이디: ${searchParams.get("orderId")}`}</div>
      <div>{`결제 금액: ${Number(
        searchParams.get("amount")
      ).toLocaleString()}원`}</div>
    </div>
  );
};

export default SuccessPaymentPage;
