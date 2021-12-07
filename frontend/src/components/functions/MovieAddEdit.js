import React, {Fragment, useEffect, useState} from "react";
import {Link, useNavigate, useParams} from "react-router-dom";
import {confirmAlert} from "react-confirm-alert";

import Input from "../form-components/Input";
import TextArea from "../form-components/TextArea";
import Select from "../form-components/Select";
import Alert from "../ui-components/Alert";

import "./css/MovieAddEdit.css"
import "react-confirm-alert/src/react-confirm-alert.css";

export default function MovieAddEdit({jwt, forbidden}) {
    let params = useParams()
    let navigate = useNavigate()

    const [movie, setMovie] = useState({})
    const [error, setError] = useState(null)
    const [errors, setErrors] = useState([])
    const [alert, setAlert] = useState({type: "d-none", message: ""})

    const ratingOptions = [
        {id: "G", value: "G"},
        {id: "PG", value: "PG"},
        {id: "PG13", value: "PG13"},
        {id: "R", value: "R"},
        {id: "NC17", value: "NC17"}
    ]
    const headers = new Headers()
    headers.append("Content-Type", "application/json")
    headers.append("Authorization", "Bearer " + jwt)

    const handleChange = (evt) => {
        let value = evt.target.value
        let name = evt.target.name
        setMovie({...movie, [name]: value})
    }

    const handleSubmit = (evt) => {
        evt.preventDefault()

        let formErrors = []
        if (movie.title.length === "")
            formErrors.push("title")

        if (movie.release_date === null)
            formErrors.push("date")

        if (ratingOptions.findIndex(item => item.id === movie.rating) === -1)
            formErrors.push("rating")

        setErrors(formErrors)
        if (formErrors.length > 0)
            return false

        const requestOptions = {
            method: "POST",
            body: JSON.stringify(Object.fromEntries(new FormData(evt.target).entries())),
            headers: headers
        }

        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/admin/movie/function`, requestOptions)
            .then(response => response.json())
            .then((data) => {
                if (data.error)
                    setAlert({alert: {type: "alert-danger", message: data.error.message}})
                else
                    navigate("/admin")
            })
    }

    const confirmDelete = () => {
        confirmAlert({
            title: 'Delete movie',
            message: 'Are you sure?',
            buttons: [
                {
                    label: 'Yes',
                    onClick: () => {
                        const requestOptions = {
                            method: "POST",
                            body: JSON.stringify({id: params.id}),
                            headers: headers
                        }
                        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/admin/movie/delete`, requestOptions)
                            .then(response => response.json())
                            .then((data) => {
                                if (data.error)
                                    setAlert({alert: {type: "alert-danger", message: data.error.message}})
                                else
                                    navigate("/admin")
                            })
                    }
                },
                {
                    label: 'No'
                }
            ]
        })
    }

    const hasError = (key) => {
        return errors.indexOf(key) !== -1
    }

    useEffect(() => {
        if (jwt === "" || forbidden) {
            navigate('/login')
            return
        }

        const id = params.id
        if (id > 0) {
            fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/movie/` + params.id)
                .then((response) => {
                    if (response.status !== 200)
                        setError("Invalid response code: ", response.status)
                    else
                        setError(null)
                    return response.json()
                })
                .then((json) => {
                    json.movie.release_date = new Date(json.movie.release_date).toISOString().split("T")[0]
                    setMovie(json.movie)
                })
        }
    }, [forbidden, jwt, navigate, params.id])

    if (error !== null)
        return <div>Error: {error.message}</div>
    else {
        return (
            <Fragment>
                <h2>Add/Edit movie</h2>
                <Alert alertType={alert.type} alertMessage={alert.message}>{}</Alert>
                <hr/>
                <form onSubmit={handleSubmit}>
                    <input type={"hidden"} name={"id"} id={"id"} value={params.id} onChange={handleChange}/>
                    <Input title={"Title"} type={'text'} name={'title'} value={movie.title}
                           handleChange={handleChange}
                           className={hasError("title") ? "is-invalid" : ""}
                           errorDiv={hasError("title") ? "text-danger" : "d-none"} errorMsg={"Please enter a title"}/>
                    <Input title={"Release date"} type={'date'} name={'release_date'} value={movie.release_date}
                           handleChange={handleChange} className={hasError("date") ? "is-invalid" : ""}
                           errorDiv={hasError("date") ? "text-danger" : "d-none"}
                           errorMsg={"Please select a release date"}/>
                    <Input title={"Runtime"} type={'text'} name={'runtime'} value={movie.runtime}
                           handleChange={handleChange}/>
                    <Select title={"Rating"} name={"rating"} value={movie.rating} handleChange={handleChange}
                            options={ratingOptions} className={hasError("rating") ? "is-invalid" : ""}
                            errorDiv={hasError("rating") ? "text-danger" : "d-none"}
                            errorMsg={"Please select a rating"}/>
                    <TextArea title={"Description"} name={"description"} rows={"3"} value={movie.description}
                              handleChange={handleChange}/>
                    <hr/>
                    <button className={"btn btn-primary"}>Save</button>
                    <Link to={"/admin"} className={"btn btn-warning ms-1"}>Cancel</Link>
                    {params.id > 0 && (
                        <a href={"#!"} onClick={() => confirmDelete()} className={"btn btn-danger ms-1"}>Delete</a>
                    )}
                </form>
            </Fragment>
        )
    }
}