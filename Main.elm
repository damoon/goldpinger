-- Read more about this program in the official Elm guide:
-- https://guide.elm-lang.org/architecture/effects/http.html


module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import Json.Decode
import Time exposing (Time, second)
import Round
import Dict exposing (Dict)


main : Program Never Model Msg
main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


init : ( Model, Cmd Msg )
init =
    ( { nodes = [], measurements = Dict.fromList [], error = "" }, fetchResults )


subscriptions : Model -> Sub Msg
subscriptions model =
    Time.every (1 * second) Tick


type alias Model =
    { nodes : List Node
    , measurements : Dict String Dict String Measurement
    , error : String
    }


type alias Node =
    { hostName : String
    , hostIP : String
    , podName : String
    , podIP : String
    }


type alias Measurement =
    { delay : Int
    , timestamp : Int
    , error : String
    }


type Msg
    = Tick Time
    | NewResults (Result Http.Error Model)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Tick newTime ->
            ( model, fetchResults )

        NewResults (Ok newResults) ->
            ( { model | hosts = newResults, error = "" }, Cmd.none )

        NewResults (Err error) ->
            ( { model | error = toString error }, Cmd.none )


view : Model -> Html Msg
view model =
    div [ class "goldpinger" ]
        [ css "https://necolas.github.io/normalize.css/8.0.0/normalize.css"
        , css "styles.css"
        , printError model.error
        , Html.h1 [] [ text "Goldpinger" ]
        , viewTable model
        ]


css : String -> Html Msg
css path =
    node "link" [ rel "stylesheet", href path ] []


printError : String -> Html Msg
printError error =
    if error == "" then
        text ""
    else
        div [ class "error" ] [ text error ]


viewTable : Model -> Html Msg
viewTable model =
    Html.table [] (List.concat [ [ viewHosts model ], List.map viewRow model.nodes ])


viewHosts : Model -> Html Msg
viewHosts model =
    Html.tr []
        (List.concat
            [ [ Html.td [] [] ]
            , List.map (\h -> Html.td [] [ div [ class "to" ] [ text "to ", text h.hostName ] ]) model.nodes
            ]
        )


viewRow : Node -> List String -> Dict String Measurement -> Html Msg
viewRow node nodes measurements =
    Html.tr []
        (List.concat
            [ [ Html.td [] [ text "from ", text node.hostName ] ]
            , List.map (printColored) measurements
            ]
        )


printColored : Measurement -> Html Msg
printColored ping =
    let
        delayInMilisec =
            toFloat ping.delay / 100000

        display =
            Round.round 2 delayInMilisec
    in
        if delayInMilisec > 50 then
            Html.td [ class "high ping" ] [ text display ]
        else if delayInMilisec > 25 then
            Html.td [ class "med ping" ] [ text display ]
        else
            Html.td [ class "low ping" ] [ text display ]


fetchResults : Cmd Msg
fetchResults =
    Http.send NewResults (Http.get "./status.json" (Json.Decode.list decodeHost))


decodeModel : Json.Decode.Decoder Model
decodeModel =
    Json.Decode.dict Dict Measurement


decodeHost : Json.Decode.Decoder Node
decodeHost =
    Json.Decode.map4 Node
        (Json.Decode.field "hostName" Json.Decode.string)
        (Json.Decode.field "hostIP" (Json.Decode.string))
        (Json.Decode.field "podName" (Json.Decode.string))
        (Json.Decode.field "podIP" (Json.Decode.string))


decodePing : Json.Decode.Decoder Measurement
decodePing =
    Json.Decode.map3 Measurement
        (Json.Decode.field "delay" Json.Decode.int)
        (Json.Decode.field "timestamp" Json.Decode.int)
        (Json.Decode.field "error" Json.Decode.string)
