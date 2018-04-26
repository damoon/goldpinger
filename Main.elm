-- Read more about this program in the official Elm guide:
-- https://guide.elm-lang.org/architecture/effects/http.html


module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import Json.Decode
import Time exposing (Time, second)


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
    ( { hosts = [], error = "" }
    , fetchResults
    )


subscriptions : Model -> Sub Msg
subscriptions model =
    Time.every (1 * second) Tick


type alias Model =
    { hosts : List Host
    , error : String
    }


type alias Host =
    { source : String
    , pings : List Ping
    }


type alias Ping =
    { target : String
    , delay : Int
    , timestamp : Int
    }


type Msg
    = Tick Time
    | NewResults (Result Http.Error (List Host))


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
        , viewTable model.hosts
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


viewTable : List Host -> Html Msg
viewTable hosts =
    Html.table [] (List.concat [ [ viewHosts hosts ], List.map viewRow hosts ])


viewHosts : List Host -> Html Msg
viewHosts hosts =
    Html.tr []
        (List.concat
            [ [ Html.td [] [] ]
            , List.map (\h -> Html.td [] [ div [ class "to" ] [ text "to ", text h.source ] ]) hosts
            ]
        )


viewRow : Host -> Html Msg
viewRow host =
    Html.tr []
        (List.concat
            [ [ Html.td [] [ text "from ", text host.source ] ]
            , List.map (viewPing host.source) host.pings
            ]
        )


viewPing : String -> Ping -> Html Msg
viewPing source ping =
    if ping.target == source then
        Html.td [ Html.Attributes.class "local-machine" ] [ text "-" ]
    else
        printColored ping.delay


printColored : Int -> Html Msg
printColored delay =
    if delay > 50 then
        Html.td [ class "high ping" ] [ text (toString delay) ]
    else if delay > 25 then
        Html.td [ class "med ping" ] [ text (toString delay) ]
    else
        Html.td [ class "low ping" ] [ text (toString delay) ]


fetchResults : Cmd Msg
fetchResults =
    let
        url =
            "http://localhost:8080/"
    in
        Http.send NewResults (Http.get url (Json.Decode.list decodeHost))


decodeHost : Json.Decode.Decoder Host
decodeHost =
    Json.Decode.map2 Host
        (Json.Decode.field "source" Json.Decode.string)
        (Json.Decode.field "pings" (Json.Decode.list decodePing))


decodePing : Json.Decode.Decoder Ping
decodePing =
    Json.Decode.map3 Ping
        (Json.Decode.field "target" Json.Decode.string)
        (Json.Decode.field "delay" Json.Decode.int)
        (Json.Decode.field "timestamp" Json.Decode.int)
