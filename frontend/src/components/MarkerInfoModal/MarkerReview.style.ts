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

const blinkAnimation = keyframes`
  0%, 50%, 100% {
    opacity: 1;
  }
  25%, 75% {
    opacity: 0.5;
  }
`;

const flexCenter = `
  display: flex;
  align-items: center;
  justify-content: center;
`;

export const Container = styled.div`
  ${flexCenter}
  flex-direction: column;

  margin-top: 3rem;
`;

export const Wrapper = styled.div`
  ${flexCenter}

  margin-bottom: 2rem;

  height: 4rem;
`;

export const P = styled.div`
  font-size: 5rem;
`;

export const DotWrap = styled.div`
  display: flex;
  align-items: end;

  height: 100%;
  width: 100%;
`;

const Dot = styled.div`
  margin: 0 0.3rem;

  width: 0.8rem;
  height: 0.8rem;

  border-radius: 50%;
  background-color: #222;

  animation: ${blinkAnimation} 2s infinite;
`;

export const Dot1 = styled(Dot)`
  animation-delay: 0s;
`;

export const Dot2 = styled(Dot)`
  animation-delay: 0.2s;
`;

export const Dot3 = styled(Dot)`
  animation-delay: 0.4s;
`;

export const Text = styled.div`
  ${flexCenter}
  flex-direction: column;

  margin-bottom: 1rem;

  & > p:first-of-type {
    margin-bottom: 0.3rem;

    font-size: 1.3rem;
    font-weight: 500;
  }
`;

export const ReviewWrap = styled.div`
  margin-bottom: 2rem;
  padding: 1rem;

  height: 300px;

  overflow: auto;
`;

export const ReviewItem = styled.div`
  display: flex;
  align-items: center;

  padding: 1rem;
  margin: 0 auto 1rem auto;

  width: 100%;

  white-space: nowrap;
  overflow: hidden;

  text-overflow: ellipsis;

  border-radius: 0.4rem;
  background-color: #e9efff;

  &:hover > div:first-of-type {
    word-wrap: break-word;
    white-space: -moz-pre-wrap;
    white-space: pre-wrap;

    text-overflow: none;
  }

  & > div:first-of-type {
    width: 60%;

    text-align: left;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
    font-size: 1.1rem;
  }
  & > div:nth-of-type(2) {
    flex-grow: 1;
  }
  & > div:nth-of-type(3) {
    font-size: 0.7rem;
    color: #777;
  }
  & > div:last-of-type {
  }
`;

export const InputWrap = styled.div`
  display: flex;
  align-items: center;

  border-radius: 1rem;
  padding: 0 0.5rem;
  border: 1.5px solid #888;
  border-radius: 0.5rem;
`;

export const ReviewInput = styled.input`
  flex-grow: 1;

  border: none;
  outline: none;

  height: 1.5rem;

  font: inherit;
`;

export const ListSkeleton = styled.div`
  display: flex;
  align-items: center;

  margin: 0 auto 1rem auto;

  padding: 1rem;

  height: 57px;
  width: 100%;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;
  border-radius: 0.4rem;
`;

export const ErrorBox = styled.div`
  text-align: left;

  font-size: 0.7rem;

  color: red;
`;
