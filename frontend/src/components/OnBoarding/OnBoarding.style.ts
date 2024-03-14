import styled from "@emotion/styled";
import { keyframes } from "@emotion/react";

const toDown = keyframes`
  0% {
    top: -100%;
  }
  100% {
    top: 5%;
  }
`;

const transparent = keyframes`
  0% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
`;

export const Title = styled.div`
  user-select: none;
  color: #fff;

  font-size: 2rem;
  font-weight: bold;

  margin-top: 12rem;

  transition: all 1s;

  & > p:first-of-type {
  }
  & > p:last-of-type {
    font-size: 1.2rem;

    margin: 2rem;
  }
`;

export const Container = styled.div`
  position: absolute;
  top: 0;
  left: 0;

  width: 100vw;
  height: 100vh;

  background-color: rgba(0, 0, 0, 0.5);

  z-index: 1000;

  & img {
    user-select: none;
  }
`;

export const Step1 = styled.div`
  position: absolute;
  top: 40%;
  left: 50%;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;
`;

export const ArrowL2 = styled.img`
  width: 100px;

  @media (max-width: 670px) {
    display: none;
  }
`;

export const R1 = styled.div`
  position: absolute;
  top: 0;
  left: 110%;

  width: 210px;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }

  @media (max-width: 670px) {
    left: -80px;
    top: 50px;
  }
`;

export const MarkerImageWrap = styled.div`
  position: absolute;
  left: -40%;

  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.7);

  animation: ${toDown} 0.5s ease-in-out 1 forwards;
`;

export const Step2 = styled.div`
  position: absolute;
  bottom: 90px;
  left: 43%;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  @media (max-width: 420px) {
    left: 23%;
  }
`;

export const ArrowCd = styled.img`
  height: 100px;
`;

export const R2 = styled.div`
  position: absolute;
  top: 0;
  left: 110%;

  width: 160px;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step3 = styled.div`
  position: absolute;
  top: 27px;
  right: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    width: 100px;
  }

  @media (max-width: 380px) {
    top: 90px;
  }
`;

export const R3 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  width: 200px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step4 = styled.div`
  position: absolute;
  top: 205px;
  right: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    width: 100px;
  }
`;

export const R4 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  width: 200px;

  border-radius: 10px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step5 = styled.div`
  position: absolute;
  top: 248px;
  right: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    width: 100px;
  }
`;

export const R5 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  width: 200px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step6 = styled.div`
  position: absolute;
  top: 310px;
  right: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    width: 100px;
  }
`;

export const R6 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  width: 200px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step7 = styled.div`
  position: absolute;
  top: 350px;
  right: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    width: 100px;
  }
`;

export const R7 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  width: 200px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step8 = styled.div`
  position: absolute;
  top: 70px;
  left: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    height: 60px;
  }
`;

export const R8 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  width: 200px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step9 = styled.div`
  position: absolute;
  top: 150px;
  left: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    height: 60px;
  }
`;

export const R9 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  width: 200px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step10 = styled.div`
  position: absolute;
  top: 70px;
  left: 53%;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    height: 60px;
  }

  @media (min-width: 880px) {
    left: 460px;
  }

  @media (max-width: 490px) {
    left: 260px;
  }
`;

export const R10 = styled.div`
  position: absolute;
  top: 100%;
  left: -400%;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  width: 200px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }

  @media (max-width: 490px) {
    left: -530%;
  }
`;

export const Step11 = styled.div`
  position: absolute;
  top: 165px;
  left: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    height: 60px;
  }
`;

export const R11 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  border-radius: 10px;

  width: 200px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;

export const Step12 = styled.div`
  position: absolute;
  top: 415px;
  right: 65px;

  animation: ${transparent} 0.5s ease-in-out 1 forwards;

  & img {
    width: 100px;
  }
`;

export const R12 = styled.div`
  position: absolute;
  top: 100%;
  left: -100%;

  background-color: #fff;
  padding: 1rem;

  width: 200px;

  border-radius: 10px;

  box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

  & > p:first-of-type {
    font-size: 1.2rem;
  }
  & > p:last-of-type {
    font-size: 0.8rem;
  }
`;
