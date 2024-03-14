import * as Styled from "./OnBoarding.style";
import useOnBoardingStore from "../../store/useOnBoardingStore";

const OnBoarding = () => {
  const onBoardingState = useOnBoardingStore();

  const handleNextStep = () => {
    onBoardingState.nextStep();

    if (onBoardingState.step === 12) {
      onBoardingState.close();
    }
  };

  return (
    <Styled.Container onClick={handleNextStep}>
      <Styled.Title
        style={{
          marginTop:
            onBoardingState.step === 0 || onBoardingState.step === 12
              ? "12rem"
              : "3rem",
        }}
      >
        <p>대한민국 철봉 지도에 오신것을 환영합니다!!</p>
        {onBoardingState.step === 0 && <p>클릭으로 설명 보기</p>}
        {onBoardingState.step === 12 && <p>시작하기</p>}
      </Styled.Title>

      {(onBoardingState.step === 1 || onBoardingState.step === 2) && (
        <Styled.Step1>
          <div>
            <Styled.ArrowL2 src="/images/arrowL2.png" alt="" />
          </div>
          <Styled.R1>
            <p>철봉 위치 등록</p>
            <p>지도의 원하는 위치를 클릭하여 철봉의 위치를 등록 하세요!</p>
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
            <p>로그인을 완료하여 등록 위치와 좋아요 한 위치를 관리해 보세요!</p>
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

      {onBoardingState.step === 8 && (
        <Styled.Step8>
          <div>
            <img src="/images/arrowcu.png" alt="" />
          </div>
          <Styled.R8>
            <p>검색</p>
            <p>원하는 위치를 검색하여 이동할 수 있습니다!</p>
          </Styled.R8>
        </Styled.Step8>
      )}

      {onBoardingState.step === 9 && (
        <Styled.Step9>
          <div>
            <img src="/images/arrowcu.png" alt="" />
          </div>
          <Styled.R9>
            <p>결과 확인</p>
            <p>검색 결과를 확인하고 해당하는 위치로 이동할 수 있습니다!</p>
          </Styled.R9>
        </Styled.Step9>
      )}

      {onBoardingState.step === 10 && (
        <Styled.Step10>
          <div>
            <img src="/images/arrowcu.png" alt="" />
          </div>
          <Styled.R10>
            <p>주변 검색</p>
            <p>
              현재 화면에 보이는 지도의 중앙을 기준으로 100m ~ 5km까지 철봉
              위치를 검색할 수 있습니다!
            </p>
          </Styled.R10>
        </Styled.Step10>
      )}

      {onBoardingState.step === 11 && (
        <Styled.Step11>
          <div>
            <img src="/images/arrowcu.png" alt="" />
          </div>
          <Styled.R11>
            <p>범위 지정</p>
            <p>범위를 지정하고 주변의 철봉 위치를 검색할 수 있습니다!</p>
          </Styled.R11>
        </Styled.Step11>
      )}

      {onBoardingState.step === 12 && (
        <Styled.Step12>
          <div>
            <img src="/images/arrowR1.png" alt="" />
          </div>
          <Styled.R12>
            <p>도움 말</p>
            <p>본 설명은 언제든지 다시 보실 수 있습니다!</p>
          </Styled.R12>
        </Styled.Step12>
      )}
    </Styled.Container>
  );
};

export default OnBoarding;
