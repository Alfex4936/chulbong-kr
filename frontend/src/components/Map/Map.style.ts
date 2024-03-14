import styled from "@emotion/styled";
import { keyframes } from "@emotion/react";

export const Container = styled.div`
  position: relative;
  min-width: 320px;
`;

export const MapContainer = styled.div`
  width: 100dvw;
  height: 100dvh;
`;

export const AlertContainer = styled.div`
  border: 1px solid red;

  position: absolute;
  left: 0;
  bottom: 1rem;

  width: 50%;
  height: 80px;

  background-color: #333;

  z-index: 1;
`;

export const ExitButton = styled.div`
  position: absolute;
  top: 0;
  right: 0;

  border-radius: 50%;

  width: 25px;
  height: 25px;

  &:hover {
    color: rgb(230, 103, 103);
  }
`;

export const DeleteUserButtonsWrap = styled.div`
  display: flex;

  & button {
    margin: 1rem 1rem 0 1rem;
  }
`;

export const ChangePasswordButtonsWrap = styled.div`
  display: flex;

  & button {
    margin: 1rem 0.5rem 0 0.5rem;
  }
`;

export const ErrorBox = styled.div`
  text-align: left;

  font-size: 0.7rem;

  color: red;
`;

export const AlertText = styled.div`
  & > p:first-of-type {
    font-size: 1.5rem;
    font-weight: bold;
  }

  & > p:last-of-type {
    font-size: 0.8rem;
    color: red;

    margin: 1rem 0;
  }
`;

export const LoginButtonWrap = styled.div`
  position: absolute;

  width: 40px;
  height: 20px;

  top: 20px;
  right: 20px;

  z-index: 10;

  @media (max-width: 380px) {
    top: 80px;
  }
`;

const rippleEffect = keyframes`
  from {
    width: 0;
    height: 0;
    opacity: 0.8;
  }
  to {
    width: 100px;
    height: 100px;
    opacity: 0;
  }
`;

export const UserLocationMarker = styled.div`
  width: 15px;
  height: 15px;
  border: 1.5px solid red; /* Reduced border width */
  border-radius: 50%;
  background-color: red;
  position: absolute;
  transform: translate(-50%, -50%);

  &:before {
    content: "";
    display: block;
    width: 0;
    height: 0;
    position: absolute;
    top: 50%; /* Center the ripple effect */
    left: 50%;
    transform: translate(
      -50%,
      -50%
    ); /* Ensure the center of the ripple is in the center of the circle */
    border: 1px solid red;
    border-radius: 50%;
    animation: ${rippleEffect} 2s infinite;
    opacity: 0;
  }
`;
