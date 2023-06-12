package tetrio

import (
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Genetic struct {
	fileName            string
	population          int
	mutationRate        float32
	mutationMagnitude   float32
	roundsPerGeneration int

	strategies []EvaluationStrategy
}

func NewGeneticEvolution(fileName string, population int, mutationRate float32, mutationMagnitude float32, roundsPerGeneration int) Genetic {
	parseFloat := func(str string) float32 {
		val, err := strconv.ParseFloat(str, 32)
		if err != nil {
			log.Fatal(err)
		}
		return float32(val)
	}

	strategies := []EvaluationStrategy{}

	//read weight from existed data
	weightData, err := os.ReadFile(fileName)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
	if len(weightData) > 0 {
		weights := strings.Split(string(weightData), "\n")
		for i := 0; i < len(weights)-1; i++ {
			params := strings.Split(weights[i], ", ")
			strategies = append(strategies, EvaluationStrategy{[4]float32{parseFloat(params[0]), parseFloat(params[1]), parseFloat(params[2]), parseFloat(params[3])}})
		}
	}

	//fill the remaining with randomly generated params
	for len(strategies) < population {
		strategies = append(strategies, randomStrategy())
	}

	return Genetic{fileName, population, mutationRate, mutationMagnitude, roundsPerGeneration, strategies}
}

func (training *Genetic) Train() {
	for generation := 0; true; generation++ {
		log.Println("generation", generation)
		training.simulate()
		log.Println("saving data", generation)
		training.save()
	}
}

func (training *Genetic) simulate() {
	//simulate the games
	type Rank struct {
		index   int
		fitness int32
	}

	baseShapes := [8]int32{I0Shape, J0Shape, L0Shape, J0Shape, O0Shape, T0Shape, S0Shape, Z0Shape}
	shapes := make([]int32, training.roundsPerGeneration)
	for i := range shapes {
		shapes[i] = baseShapes[rand.Intn(len(baseShapes))]
	}
	rank := make([]Rank, 0, training.population)
	for i := range training.strategies {
		log.Printf("%v | %v\n", training.population, i)
		rank = append(rank, Rank{i, training.fitness(training.strategies[i], shapes)})
	}

	//top 3 strategies survivor for another round
	sort.Slice(rank, func(i, j int) bool { return rank[i].fitness > rank[j].fitness })
	log.Printf("high score: %v\n", rank[0].fitness)
	breeds := []EvaluationStrategy{
		training.strategies[rank[0].index],
		training.strategies[rank[1].index],
		training.strategies[rank[2].index],
	}

	//distribute top 2 strategies's weight params to the top 25% until reaches 80% of the population, mutate by chance
	canBreedPercentile := int(float32(training.population) * 0.25)
	targetBreedPercentile := int(float32(training.population) * 0.80)

	for len(breeds) < targetBreedPercentile {
		breed := training.strategies[rank[rand.Intn(canBreedPercentile)].index]
		crossGen := rand.Intn(len(breed.weights))
		diff := (training.strategies[rank[rand.Intn(2)].index].weights[crossGen] - breed.weights[crossGen]) * 0.05
		breed.weights[crossGen] += diff
		if rand.Float32() < training.mutationRate {
			diffMagnitude := breed.weights[crossGen] * (rand.Float32()*2*training.mutationMagnitude - training.mutationMagnitude)
			breed.weights[crossGen] += diffMagnitude
		}
		breeds = append(breeds, breed)
	}

	//add random mutation to fill the remaining population
	for len(breeds) < training.population {
		breeds = append(breeds, randomStrategy())
	}

	training.strategies = breeds
}

func (training *Genetic) fitness(strategy EvaluationStrategy, shapes []int32) int32 {
	const kDepth = 2
	tetris := Tetris{EvaluationStrategy: strategy}

	round := int32(0)
	for len(shapes) >= kDepth {
		result := tetris.FindMove(shapes[:kDepth])
		if !result.IsDead() {
			tetris.MakeMove(result.Shape, result.Column)
			shapes = shapes[1:]
			round++
		} else {
			break
		}
	}
	return round
}

func (training *Genetic) save() {
	file, err := os.Create(training.fileName)
	if err != nil {
		log.Fatal(err)
	}
	for i := range training.strategies {
		_, err := file.WriteString(training.strategies[i].String() + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func randomStrategy() EvaluationStrategy {
	randomParam := func() float32 {
		return rand.Float32()*2.0 - 1.0
	}
	strategy := NewEvaluationStrategy(randomParam(), randomParam(), randomParam(), randomParam())
	return strategy
}
