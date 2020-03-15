import React from 'react';
import { Container, Row, Col, Form, Button, Alert } from "react-bootstrap";
import Ajax from '../Ajax';

interface LoginFormState {
    isLoading: boolean
    hasDataError: boolean
    email: string
    otp: string
    requireOtp: boolean
    password: string
}

export default class LoginForm extends React.Component<{}, LoginFormState> {
    constructor(props: any) {
        super(props);
        this.state = {
            isLoading: false,
            hasDataError: false,
            email: '',
            otp: '',
            requireOtp: false,
            password: ''
        };
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event: React.FormEvent) {
        let target = event.target as HTMLInputElement;
        this.setState({[target.name]: target.value} as React.ComponentState );
        event.preventDefault();
    }

    handleSubmit(event: React.FormEvent) {
        event.preventDefault();
        this.setState({
            isLoading: true,
            hasDataError: false
        });
        let data = {
            email: this.state.email,
            password: this.state.password,
            otp: this.state.otp
        };
        Ajax.postData("/auth/login", data).then(res => {
            if (res.status === 200) {
                if (res.json.otpRequired) {
                    this.setState({
                        requireOtp: true,
                        isLoading: false
                    });
                } else {
                    window.sessionStorage.setItem("accessToken", res.json.accessToken);
                    window.sessionStorage.setItem("refreshToken", res.json.refreshToken);
                    window.location.href = "/dashboard.html";
                }
                return;
            }
            this.setState({
                hasDataError: true,
                isLoading: false
            });
        }).catch((e) => {
            this.setState({
                hasDataError: true,
                isLoading: false
            });
        });
    }

    render() {
        return (
            <Container>
                <Row className="justify-content-md-center">
                    <Col lg="5">
                        <h1 className="display-4 font-weight-normal">Log In</h1>
                        <Alert variant="danger" show={this.state.hasDataError}>Please verify the data you have entered.</Alert>
                        <Form onSubmit={this.handleSubmit}>
                            <Form.Group controlId="email" hidden={this.state.requireOtp}>
                                <Form.Label>Email address</Form.Label>
                                <Form.Control type="email" name="email" placeholder="your@email.address" value={this.state.email} onChange={this.handleChange} autoFocus={true} required={true} />
                            </Form.Group>
                            <Form.Group controlId="password" hidden={this.state.requireOtp}>
                                <Form.Label>Password</Form.Label>
                                <Form.Control type="password" name="password" placeholder="Enter your password" value={this.state.password} onChange={this.handleChange} minLength={8} maxLength={32} required={true} />
                            </Form.Group>
                            <Form.Group controlId="otp" hidden={!this.state.requireOtp}>
                                <Form.Label>Six-Digit Code</Form.Label>
                                <Form.Control type="text" pattern="\d*" name="otp" placeholder="Time-based One-time Password (TOTP)" value={this.state.otp} onChange={this.handleChange} minLength={6} maxLength={6} required={this.state.requireOtp} />
                            </Form.Group>
                            <Button variant="primary" type="submit" disabled={this.state.isLoading}>{this.state.isLoading ? 'Loading...' : 'Submit'}</Button>
                        </Form>
                    </Col>
                </Row>
            </Container>
        );
    }
}
