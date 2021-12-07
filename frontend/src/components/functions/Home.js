import React, {Fragment, useState} from "react";

import Select from "../form-components/Select";
import Input from "../form-components/Input";

export default function Home() {
    const [channel, setChannel] = useState({channelName: "", maxResults: ""})
    const [videos, setVideos] = useState([])
    const [error, setError] = useState(null)
    const [errors, setErrors] = useState([])

    const channels = [
        {id: "Filme", value: "KinoCheck.com"},
        {id: "movieclipsTRAILERS", value: "Movieclips Trailers"},
        {id: "SonyPictures", value: "Sony Pictures Entertainment"}
    ]

    const hasError = (key) => {
        return errors.indexOf(key) !== -1
    }

    const handleChange = (evt) => {
        let value = evt.target.value
        let name = evt.target.name
        setChannel({...channel, [name]: value})
    }

    const handleSubmit = (evt) => {
        evt.preventDefault()

        let formErrors = []
        if (channel.channelName === "")
            formErrors.push("channelName")

        if (channel.maxResults === "")
            formErrors.push("maxResults")

        setErrors(formErrors)
        if (formErrors.length > 0)
            return false

        fetch(`${process.env.REACT_APP_YT_API}/channels?key=${process.env.REACT_APP_YT_API_KEY}&part=snippet%2CcontentDetails%2Cstatistics&forUsername=${channel.channelName}`)
            .then((response) => {
                if (response.status !== 200)
                    setError("Invalid response code: ", response.status)
                else
                    setError(null)

                return response.json()
            })
            .then((json) => {
                fetch(`${process.env.REACT_APP_YT_API}/search?key=${process.env.REACT_APP_YT_API_KEY}&channelId=${json.items[0].id}&part=snippet,id&order=date&maxResults=${channel.maxResults}`)
                    .then((response) => {
                        if (response.status !== 200)
                            setError("Invalid response code: ", response.status)
                        else
                            setError(null)

                        return response.json()
                    })
                    .then((json) => {
                        setVideos(json.items)
                    })
            })
    }

    if (error !== null)
        return <div>Error: {error.message}</div>
    else {
        return (
            <Fragment>
                <h2>New movie trailers</h2>
                <hr/>
                <div className={"container"}>
                    <form onSubmit={handleSubmit}>
                        <div className={"row"}>
                            <div className={"col"}>
                                <Select title={"Select a channel"} name={"channelName"} value={channel.channelName}
                                        handleChange={handleChange} options={channels}
                                        className={hasError("channelName") ? "is-invalid" : ""}
                                        errorDiv={hasError("channelName") ? "text-danger" : "d-none"}
                                        errorMsg={"Please select a channel name"}/>
                            </div>
                            <div className={"col"}>
                                <Input title={"Max results"} type={"number"} name={"maxResults"}
                                       value={channel.maxResults} handleChange={handleChange}
                                       className={hasError("maxResults") ? "is-invalid" : ""}
                                       errorDiv={hasError("maxResults") ? "text-danger" : "d-none"}
                                       errorMsg={"Please write how many results do you want to see"}/>
                            </div>
                            <div className={"col"}>
                                <div className={"mt-2"}>
                                    <br/>
                                    <button className={"btn btn-primary"}>Search</button>
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
                <div className={"list-group"}>
                    {videos.map((v) => (
                        <div className={"list-group-item"} key={v.id.videoId}>
                            <div className={"container"}>
                                <div className={"row"}>
                                    <div className={"col-sm-8"}>
                                        <strong>{v.snippet.title}</strong><br/>
                                        <small>Published: {new Date(v.snippet.publishTime).toISOString().split("T")[0]}</small><br/>
                                        {v.snippet.description}
                                    </div>
                                    <div className={"col-sm-4"}>
                                        <img src={v.snippet.thumbnails.default.url} alt={""}/><br/>
                                        <a href={`https://www.youtube.com/watch?v=${v.id.videoId}`} target={"_blank"}>Watch
                                            the trailer</a>
                                    </div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </Fragment>
        )
    }
}