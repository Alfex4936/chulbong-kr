import * as Styled from "./SupportModal.style";
import ActionButton from "../ActionButton/ActionButton";
import { useEffect, useRef, useState } from "react";
import {
  PaymentWidgetInstance,
  loadPaymentWidget,
} from "@tosspayments/payment-widget-sdk";
import { nanoid } from "nanoid";
import useGetMyInfo from "../../hooks/query/user/useGetMyInfo";

const selector = "#payment-widget";
const clientKey = "test_gck_docs_Ovk5rk1EwkEbP0W43n07xlzm";
const customerKey = "YbX2HuSlsC9uVJW6NMRMj";

const SupportModal = () => {
  const { data, isLoading } = useGetMyInfo();

  const paymentWidgetRef = useRef<PaymentWidgetInstance | null>(null);
  const paymentMethodsWidgetRef = useRef<ReturnType<
    PaymentWidgetInstance["renderPaymentMethods"]
  > | null>(null);

  const [curPrice, setCurPrice] = useState<number>(1000);

  useEffect(() => {
    (async () => {
      const paymentWidget = await loadPaymentWidget(clientKey, customerKey);

      const paymentMethodsWidget = paymentWidget.renderPaymentMethods(
        selector,
        { value: curPrice, currency: "KRW", country: "KR" },
        { variantKey: "DEFAULT" }
      );

      paymentWidgetRef.current = paymentWidget;
      paymentMethodsWidgetRef.current = paymentMethodsWidget;
    })();
  }, []);

  useEffect(() => {
    const paymentMethodsWidget = paymentMethodsWidgetRef.current;

    if (paymentMethodsWidget == null) {
      return;
    }

    paymentMethodsWidget.updateAmount(curPrice);
  }, [curPrice]);

  if (isLoading) return <p>로딩중...</p>;

  return (
    <div>
      {/* <div ref={paymentWidgetRef}></div>
      <div ref={paymentMethodsWidgetRef}></div> */}
      <div id="payment-widget" />
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
      <ActionButton
        bg="black"
        onClick={async () => {
          const paymentWidget = paymentWidgetRef.current;
          console.log(paymentWidget);

          // try {
          //   // ## Q. 결제 요청 후 계속 로딩 중인 화면이 보인다면?
          //   // 아직 결제 요청 중이에요. 이어서 요청 결과를 확인한 뒤, 결제 승인 API 호출까지 해야 결제가 완료돼요.
          //   // 코드샌드박스 환경에선 요청 결과 페이지(`successUrl`, `failUrl`)로 이동할 수가 없으니 유의하세요.
          //   await paymentWidget?.requestPayment({
          //     orderId: nanoid(),
          //     orderName: "커피 한잔 후원~",
          //     customerName: data?.username,
          //     customerEmail: data?.email,
          //     // successUrl: `${window.location.origin}/success`,
          //     // failUrl: `${window.location.origin}/fail`,
          //   });
          // } catch (error) {
          //   // handle error
          // }
        }}
      >
        {curPrice}원 후원하기
      </ActionButton>
    </div>
  );
};

export default SupportModal;
