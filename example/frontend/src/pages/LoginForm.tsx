import React from 'react';
import { Container, Row, Col, Form, Button, Alert } from "react-bootstrap";
import Ajax from '../Ajax';

interface LoginFormState {
    isLoading: boolean
    hasDataError: boolean
    email: string
    password: string
}

export default class LoginForm extends React.Component<{}, LoginFormState> {
    constructor(props: any) {
        super(props);
        this.state = {
            isLoading: false,
            hasDataError: false,
            email: '',
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
            password: this.state.password
        };
        Ajax.postData("/auth/login", data).then(res => {
            if (res.status === 200) {
                window.sessionStorage.setItem("accessToken", res.json.accessToken);
                window.sessionStorage.setItem("refreshToken", res.json.refreshToken);
                window.location.href = "/dashboard.html";
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
                            <Form.Group controlId="email">
                                <Form.Label>Email address</Form.Label>
                                <Form.Control type="email" name="email" placeholder="your@email.address" value={this.state.email} onChange={this.handleChange} autoFocus={true} />
                            </Form.Group>
                            <Form.Group controlId="password">
                                <Form.Label>Password</Form.Label>
                                <Form.Control type="password" name="password" placeholder="Choose a password" value={this.state.password} onChange={this.handleChange} minLength={8} maxLength={32} />
                            </Form.Group>
                            <Button variant="primary" type="submit" disabled={this.state.isLoading}>{this.state.isLoading ? 'Loading...' : 'Submit'}</Button>
                        </Form>
                    </Col>
                </Row>
            </Container>
        );
    }
}
