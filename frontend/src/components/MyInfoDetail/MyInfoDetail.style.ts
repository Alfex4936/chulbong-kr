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
  & > div:last-of-type {
    border-bottom: none;
  }
`;

export const NameContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;

  padding: 1rem;

  border-bottom: 2px solid #eee;

  & > div:first-of-type {
    font-size: 1.2rem;

    margin-right: 0.3rem;
  }
`;

export const Name = styled.div`
  white-space: nowrap;
  overflow: hidden;

  max-width: 180px;

  text-overflow: ellipsis;

  font-weight: 700;

  & > span {
    font-size: 1rem;

    font-weight: 400;
  }
`;

export const NameButtonContainer = styled.div`
  display: flex;
  justify-content: center;

  & button {
    width: 50px;

    margin: 0.5rem 0.5rem 0 0.5rem;

    font-size: 0.7rem;
  }
`;

export const EmailContainer = styled.div`
  padding: 1rem;

  border-bottom: 2px solid #eee;

  & > div:last-child {
    font-weight: 700;
  }
`;

export const PaymentContainer = styled.div`
  padding: 1rem;

  border-bottom: 2px solid #eee;
`;

export const ButtonContainer = styled.div`
  padding: 1rem 2rem;

  border-bottom: 2px solid #eee;
`;

export const ButtonTop = styled.div`
  display: flex;
  justify-content: space-between;

  & > button {
    margin: 0 0 0.8rem 0;

    width: 47%;

    font-size: 0.8rem;
  }
`;

export const ButtonBottom = styled.div``;

export const ErrorBox = styled.div`
  text-align: left;

  font-size: 0.7rem;

  padding: 0 1rem;

  color: red;
`;

export const InfoContainer = styled.div`
  padding: 1rem;

  & > p {
    font-size: 0.8rem;
    color: #666;
  }
`;

export const ListSkeleton = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;

  padding: 1rem;

  height: 60px;
  width: 100%;

  border-bottom: 1px solid #ccc;

  & > div:first-of-type {
    height: 24px;
    width: 90px;

    margin-right: 1rem;

    background: #f6f7f8;
    background-image: linear-gradient(
      90deg,
      #f0f0f0 25%,
      #f7f7f7 50%,
      #f0f0f0 75%
    );

    animation: ${shimmer} 1.2s ease-in-out infinite;

    border-radius: 1rem;
  }

  & > div:last-of-type {
    height: 25px;
    width: 25px;

    background: #f6f7f8;
    background-image: linear-gradient(
      90deg,
      #f0f0f0 25%,
      #f7f7f7 50%,
      #f0f0f0 75%
    );

    animation: ${shimmer} 1.2s ease-in-out infinite;

    border-radius: 50%;
  }
`;
