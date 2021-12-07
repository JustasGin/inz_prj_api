import React, {Fragment, useState} from "react";
import {useNavigate} from "react-router-dom";

import Alert from "../ui-components/Alert";
import Input from "../form-components/Input";

export default function Login({updateCookie}) {
    const navigate = useNavigate()

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [errors, setErrors] = useState([])
    const [alert, setAlert] = useState({type: "d-none", message: ""})

    const handleEmail = (evt) => {
        setEmail(evt.target.value)
    }

    const handlePassword = (evt) => {
        setPassword(evt.target.value)
    }

    const handleSubmit = (evt) => {
        evt.preventDefault()

        let formErrors = []
        if (email === "")
            formErrors.push("email")
        if (password === "")
            formErrors.push("password")
        setErrors(formErrors)

        if (formErrors.length > 0)
            return false

        const requestOptions = {
            method: "POST",
            body: JSON.stringify(Object.fromEntries(new FormData(evt.target).entries()))
        }

        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/signin`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    setAlert({type: "alert-danger", message: data.error.message})
                } else {
                    updateCookie(Object.values(data)[0].toString())
                    navigate("/")
                }
            })
    }

    let hasError = (key) => {
        return errors.indexOf(key) !== -1
    }

    return (
        <Fragment>
            <h2>Login</h2>
            <hr/>
            <Alert alertType={alert.type} alertMessage={alert.message}/>
            <form className={"pt-3"} onSubmit={handleSubmit}>
                <Input title={"Email"} type={"email"} name={"email"} handleChange={handleEmail}
                       className={hasError("email") ? "is-invalid" : ""}
                       errorDiv={hasError("email") ? "text-danger" : "d-none"}
                       errorMsg={"Please enter a valid email address"}/>
                <Input title={"Password"} type={"password"} name={"password"}
                       handleChange={handlePassword}
                       className={hasError("password") ? "is-invalid" : ""}
                       errorDiv={hasError("password") ? "text-danger" : "d-none"}
                       errorMsg={"Please enter a password"}/>
                <hr/>
                <button className={"btn btn-primary"}>Login</button>
            </form>
        </Fragment>
    )
}