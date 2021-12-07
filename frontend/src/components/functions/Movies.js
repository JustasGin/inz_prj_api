import React, {Fragment, useEffect, useState} from "react";
import {Link} from "react-router-dom";

export default function Movies() {
    const [movies, setMovies] = useState([])
    const [error, setError] = useState(null)

    useEffect(() => {
        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/movies`)
            .then((response) => {
                if (response.status !== 200)
                    setError("Invalid response code: ", response.status)
                else
                    setError(null)

                return response.json()
            })
            .then((json) => {
                setMovies(json.movies)
            })
    }, [])

    if (error !== null)
        return <div>Error: {error.message}</div>
    else {
        return (
            <Fragment>
                <h2>Choose a Movie</h2>
                <div className={"list-group"}>
                    {movies.map((m) => (
                        <Link to={`/movies/${m.id}`} className={"list-group-item list-group-item-action"} key={m.id}>
                            <strong>{m.title}</strong><br/>
                            <small className={"text-muted"}>
                                ({m.year}) - {m.runtime} minutes
                            </small><br/>
                            {m.description.slice(0, 100)}...
                        </Link>
                    ))}
                </div>
            </Fragment>
        )
    }
}