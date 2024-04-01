import styled from "@emotion/styled";

export const Container = styled.div`
  position: absolute;
  top: 20px;
  left: 550px;

  display: flex;
  align-items: center;
  justify-content: center;

  height: 40px;

  background-color: #fff;

  padding-inline: 1rem;
  border-radius: 0.5rem;
  z-index: 200;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px,
    rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
`;

export const WeatherWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;

  & div {
    font-size: 0.9rem;
  }

  & span {
    display: flex;
    align-items: center;
    justify-content: center;
  }
`;