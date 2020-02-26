import React from 'react';
import { Container, Row, Col, Form, Button, Alert } from "react-bootstrap";
import Ajax from '../Ajax';

interface SignupFormState {
    isLoading: boolean
    hasDataError: boolean
    hasAccountExists: boolean
    hasSuccess: boolean
    email: string
    password: string
}

export default class SignupForm extends React.Component<{}, SignupFormState> {
    constructor(props: any) {
        super(props);
        this.state = {
            isLoading: false,
            hasDataError: false,
            hasAccountExists: false,
            hasSuccess: false,
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
            hasDataError: false,
            hasAccountExists: false,
            hasSuccess: false
        });
        let data = {
            email: this.state.email,
            password: this.state.password
        };
        Ajax.postData("/auth/signup", data).then(res => {
            if (res.status === 201) {
                this.setState({hasSuccess: true});
                return;
            }
            if (res.status === 409) {
                this.setState({
                    hasAccountExists: true,
                    isLoading: false
                });
                return;
            }
            this.setState({
                hasDataError: true,
                isLoading: false
            });
        }).catch(() => {
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
                        <h1 className="display-4 font-weight-normal">Sign Up</h1>
                        <Alert variant="success" show={this.state.hasSuccess}>Please check your emails to activate your account.</Alert>
                        <Alert variant="danger" show={this.state.hasAccountExists}>Account already exists for this email address.</Alert>
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
                            { this.state.hasSuccess ? null : <Button variant="primary" type="submit" disabled={this.state.isLoading}>{this.state.isLoading ? 'Loading...' : 'Submit'}</Button> }
                        </Form>
                    </Col>
                </Row>
            </Container>
        );
    }
}
