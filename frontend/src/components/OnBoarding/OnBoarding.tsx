import * as Styled from "./OnBoarding.style";
import useOnBoardingStore from "../../store/useOnBoardingStore";

const OnBoarding = () => {
  const onBoardingState = useOnBoardingStore();

  const handleNextStep = () => {
    onBoardingState.nextStep();
  };

  return (
    <Styled.Container onClick={handleNextStep}>
      <h1>온보딩</h1>
      {(onBoardingState.step === 1 || onBoardingState.step === 2) && (
        <Styled.Step1>
          <div>
            <Styled.ArrowL2 src="/images/arrowL2.png" alt="" />
          </div>
          <Styled.R1>
            <p>지도 클릭</p>
            <p>맵의 원하는 위치를 클릭하여 위치를 등록 하세요!</p>
          </Styled.R1>
          <Styled.MarkerImageWrap>
            <img src="/images/cb2.webp" alt="" width={40} />
          </Styled.MarkerImageWrap>
          <div />
        </Styled.Step1>
      )}
      {onBoardingState.step === 2 && (
        <Styled.Step2>
          <div>
            <Styled.ArrowCd src="/images/arrowcd.png" alt="" />
          </div>
          <Styled.R2>
            <p>클릭</p>
            <p>로그인이 필요합니다!</p>
          </Styled.R2>
        </Styled.Step2>
      )}

      {onBoardingState.step === 3 && (
        <Styled.Step3>
          <div>
            <img src="/images/arrowR1.png" alt="" />
          </div>
          <Styled.R3>
            <p>로그인</p>
            <p>로그인을 완료하면 내 정보 메뉴를 볼 수 있습니다!</p>
          </Styled.R3>
        </Styled.Step3>
      )}

      {onBoardingState.step === 4 && (
        <Styled.Step4>
          <div>
            <img src="/images/arrowR1.png" alt="" />
          </div>
          <Styled.R4>
            <p>내 위치</p>
            <p>현재 내 위치를 추적하여 지도를 확인할 수 있습니다!</p>
          </Styled.R4>
        </Styled.Step4>
      )}

      {onBoardingState.step === 5 && (
        <Styled.Step5>
          <div>
            <img src="/images/arrowR1.png" alt="" />
          </div>
          <Styled.R5>
            <p>위치 초기화</p>
            <p>
              화면의 보이는 위치를 서울특별시청이 중앙에 오도록 초기화합니다!
            </p>
          </Styled.R5>
        </Styled.Step5>
      )}

      {onBoardingState.step === 6 && (
        <Styled.Step6>
          <div>
            <img src="/images/arrowR1.png" alt="" />
          </div>
          <Styled.R6>
            <p>확대</p>
            <p>지도를 확대할 수 있습니다!</p>
          </Styled.R6>
        </Styled.Step6>
      )}

      {onBoardingState.step === 7 && (
        <Styled.Step7>
          <div>
            <img src="/images/arrowR1.png" alt="" />
          </div>
          <Styled.R7>
            <p>축소</p>
            <p>지도를 축소할 수 있습니다!</p>
          </Styled.R7>
        </Styled.Step7>
      )}
    </Styled.Container>
  );
};

export default OnBoarding;
