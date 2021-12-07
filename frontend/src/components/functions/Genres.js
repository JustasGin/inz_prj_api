import React, {Fragment, useEffect, useState} from "react";
import {Link} from "react-router-dom";

export default function Genres() {
    const [genres, setGenres] = useState([])
    const [error, setError] = useState(null)

    useEffect(() => {
        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/genres`)
            .then((response) => {
                if (response.status !== 200)
                    setError("Invalid response code: ", response.status)
                else
                    setError(null)
                return response.json()
            })
            .then((json) => {
                setGenres(json.genres)
            })
    }, [])

    if (error !== null)
        return <div>Error: {error.message}</div>
    else {
        return (
            <Fragment>
                <h2>Genres</h2>
                <div className={"list-group"}>
                    {genres.map((m, index) => (
                        <Link to={`/genres/${m.id}`} state={{genreName: m.name}}
                              className={"list-group-item list-group-item-action"} key={index}>{m.name}</Link>
                    ))}
                </div>
            </Fragment>
        )
    }
}