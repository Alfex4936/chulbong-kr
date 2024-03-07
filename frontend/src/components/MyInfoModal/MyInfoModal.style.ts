import styled from "@emotion/styled";
import { keyframes } from "@emotion/react";

const shimmer = keyframes`
  0% {
    background-position: -468px 0;
  }
  100% {
    background-position: 468px 0;
  }
`;

export const Container = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;

  position: absolute;

  width: 300px;

  top: 25px;
  right: 70px;

  background-color: #fff;

  border: 1px solid #ddd;
  border-radius: 9px;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px,
    rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;

  z-index: 10;
`;

export const InfoTop = styled.div`
  display: flex;
  align-items: center;

  border-bottom: 1px solid #eee;

  padding: 1rem;
`;

export const ProfileImgBox = styled.div`
  border-radius: 50%;

  box-shadow: rgba(0, 0, 0, 0.02) 0px 1px 3px 0px,
    rgba(27, 31, 35, 0.15) 0px 0px 0px 1px;

  width: 40px;
  height: 40px;

  margin-right: 0.5rem;

  overflow: hidden;

  & img {
    display: inline-block;

    width: 100%;
  }
`;

export const NameContainer = styled.div`
  margin-right: 1rem;

  text-align: left;

  font-weight: bold;

  & > div:last-of-type {
    width: 100%;

    font-size: 0.7rem;
    font-weight: 400;
    color: #555;
  }

  flex-grow: 1;
`;

export const ButtonContainer = styled.div`
  & > button {
    font-size: 0.7rem;
  }
`;

export const InfoBottom = styled.div`
  display: flex;
`;

export const TabContainer = styled.div``;

export const NameSkeleton = styled.div`
  flex-grow: 1;
  & > div:first-of-type {
    height: 22px;
    width: 90px;

    margin-bottom: 0.3rem;

    background: #f6f7f8;
    background-image: linear-gradient(
      90deg,
      #f0f0f0 25%,
      #f7f7f7 50%,
      #f0f0f0 75%
    );

    animation: ${shimmer} 1.2s ease-in-out infinite;

    border-radius: 0.5rem;
  }

  & > div:last-of-type {
    height: 15px;
    width: 100px;

    background: #f6f7f8;
    background-image: linear-gradient(
      90deg,
      #f0f0f0 25%,
      #f7f7f7 50%,
      #f0f0f0 75%
    );

    animation: ${shimmer} 1.2s ease-in-out infinite;

    border-radius: 0.5rem;
  }
`;
