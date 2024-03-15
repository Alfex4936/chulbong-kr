import styled from "@emotion/styled";

export const Container = styled.div`
  position: absolute;
  top: 0;
  left: 0;

  width: 100%;
  height: 100%;

  backgroud-color: rgba(0, 0, 0, 0.5);

  z-index: 900;
`;

export const RoadViewContainer = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);

  width: 800px;
  height: 500px;

  @media (max-width: 900px) {
    width: 600px;
    height: 400px;
  }
  @media (max-width: 650px) {
    width: 500px;
    height: 350px;
  }
  @media (max-width: 530px) {
    width: 400px;
    height: 300px;
  }
  @media (max-width: 410px) {
    width: 320px;
    height: 320px;
  }
  @media (max-width: 330px) {
    width: 250px;
    height: 250px;
  }
`;

export const Exit = styled.div`
  position: absolute;
  bottom: -20rem;
  left: 50%;
  transform: translateX(-50%);

  @media (max-width: 900px) {
    bottom: -16rem;
  }
  @media (max-width: 650px) {
    bottom: -15rem;
  }
  @media (max-width: 530px) {
    bottom: -13rem;
  }
  @media (max-width: 330px) {
    bottom: -12rem;
  }
`;
