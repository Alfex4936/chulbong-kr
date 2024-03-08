import * as Styled from "./SupportPage.style.ts";
import ActionButton from "../ActionButton/ActionButton";
import { useState } from "react";
import { loadTossPayments } from "@tosspayments/payment-sdk";
import { nanoid } from "nanoid";
import useGetMyInfo from "../../hooks/query/user/useGetMyInfo";

const clientKey = "test_gck_docs_Ovk5rk1EwkEbP0W43n07xlzm";

const SupportPage = () => {
  const { data, isLoading } = useGetMyInfo();

  const [curPrice, setCurPrice] = useState<number>(1000);

  const handleClick = async () => {
    if (!data) return;

    const tossPayments = await loadTossPayments(clientKey);

    await tossPayments.requestPayment("카드", {
      amount: curPrice,
      orderId: nanoid(),
      orderName: "커피 한잔 후원~",
      customerName: data?.username,
      customerEmail: data?.email,
      successUrl: `${window.location.origin}/api/payments`,
      failUrl: `${window.location.origin}/api/payments/fail`,
    });
  };

  if (isLoading) return <p>로딩중...</p>;

  return (
    <div>
      <Styled.ImageBox>
        <button
          onClick={() => {
            setCurPrice(1000);
          }}
        >
          {curPrice === 1000 ? (
            <img src="/images/1000a.webp" alt="1000원" />
          ) : (
            <img src="/images/1000.webp" alt="1000원" />
          )}
        </button>
        <button
          onClick={() => {
            setCurPrice(3000);
          }}
        >
          {curPrice === 3000 ? (
            <img src="/images/3000a.webp" alt="3000원" />
          ) : (
            <img src="/images/3000.webp" alt="3000원" />
          )}
        </button>
        <button
          onClick={() => {
            setCurPrice(5000);
          }}
        >
          {curPrice === 5000 ? (
            <img src="/images/5000a.webp" alt="5000원" />
          ) : (
            <img src="/images/5000.webp" alt="5000원" />
          )}
        </button>
      </Styled.ImageBox>
      <ActionButton bg="black" onClick={handleClick}>
        후원하기
      </ActionButton>
    </div>
  );
};

export default SupportPage;
