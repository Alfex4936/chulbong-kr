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
  max-height: 400px;
  overflow: auto;
`;

export const ButtonContainer = styled.div`
  display: flex;
  justify-content: center;

  margin-bottom: 2rem;

  & > button {
    width: 200px;
    margin: 0.5rem;
  }
`;

export const MessageRed = styled.p`
  font-size: 0.8rem;
  color: #ff6060;

  padding: 0 1rem;
`;

export const ResultItem = styled.div`
  display: flex;
  align-items: center;

  padding: 1rem;
  margin: 0 auto 1rem auto;

  width: 90%;

  border-radius: 0.4rem;
  background-color: #e9efff;

  & > span:nth-of-type(2) {
    flex-grow: 1;

    font-size: 0.9rem;
    color: #777;
  }

  @media (max-width: 660px) {
    & > span:first-of-type {
      font-size: 0.8rem;
    }
    & > span:nth-of-type(2) {
      flex-grow: 1;

      font-size: 0.7rem;
      color: #777;
    }
  }
`;

export const ListSkeleton = styled.div`
  display: flex;
  align-items: center;

  margin: 0 auto 1rem auto;

  padding: 1rem;

  height: 57px;
  width: 90%;

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
