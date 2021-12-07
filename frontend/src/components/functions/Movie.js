import React, {Fragment, useEffect, useState} from "react";
import {useParams} from "react-router-dom";

export default function Movie() {
    let params = useParams()
    const [movie, setMovie] = useState({})
    const [reviews, setReviews] = useState([])
    const [error, setError] = useState(null)

    useEffect(() => {
        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/movie/${params.id}`)
            .then((response) => {
                if (response.status !== 200)
                    setError("Invalid response code: ", response.status)
                else
                    setError(null)
                return response.json()
            })
            .then((json) => {
                setMovie(json.movie)

                fetch(encodeURI(`${process.env.REACT_APP_IMDB_API}/Search/${process.env.REACT_APP_IMDB_API_KEY}/${json.movie.title}`))
                    .then((response) => {
                        if (response.status !== 200)
                            setError("Invalid response code: ", response.status)
                        else
                            setError(null)
                        return response.json()
                    })
                    .then((json) => {
                        fetch(`${process.env.REACT_APP_IMDB_API}/Reviews/${process.env.REACT_APP_IMDB_API_KEY}/${json.results[0].id}`)
                            .then((response) => {
                                if (response.status !== 200)
                                    setError("Invalid response code: ", response.status)
                                else
                                    setError(null)
                                return response.json()
                            })
                            .then((json) => {
                                setReviews(json.items)
                            })
                    })
            })
    }, [params.id])

    if (movie.genres)
        movie.genres = Object.values(movie.genres)
    else
        movie.genres = []

    if (error !== null)
        return <div>Error: {error.message}</div>
    else {
        return (
            <Fragment>
                <h2>{movie.title} ({movie.year})</h2>
                <div className={"float-end"}>
                    {movie.genres.map((m, index) => (
                        <span className={"badge bg-secondary me-1"} key={index}>{m}</span>
                    ))}
                </div>
                <small className={"text-muted"}>{movie.rating}</small>
                <hr/>
                {movie.poster !== "" && (
                    <div>
                        <img src={`https://image.tmdb.org/t/p/w200${movie.poster}`} alt={"poster"}/>
                    </div>
                )}
                <hr/>
                <table className={"table table-compact table-striped"}>
                    <tbody>
                    <tr>
                        <td><strong>Title:</strong></td>
                        <td>{movie.title}</td>
                    </tr>
                    <tr>
                        <td><strong>Description:</strong></td>
                        <td>{movie.description}</td>
                    </tr>
                    <tr>
                        <td><strong>Runtime:</strong></td>
                        <td>{movie.runtime} minutes</td>
                    </tr>
                    </tbody>
                </table>
                {reviews.map((r, index) => (
                    <div className={"card"} key={index}>
                        <div className={"card-body"}>
                            <h5 className={"card-title"}>{r.username}</h5>
                            <h6 className={"card-subtitle mb-2 text-muted"}>Rating - {r.rate}/10</h6>
                            <p className={"card-text"}>{r.content}</p>
                        </div>
                    </div>
                ))}
            </Fragment>
        )
    }
}