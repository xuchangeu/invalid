package test

//func loadTestCase(files []string) (*invalid.YAMLRoot, error) {
//
//	file, err := os.Open(filepath.Join(files...))
//	if err != nil {
//		return nil, err
//	}
//
//	r, err := invalid.NewYAML(file)
//	if err != nil {
//		return nil, err
//	}
//
//	return r, nil
//}
//
//func loadTestRule(files []string) (*invalid.RuleRoot, error) {
//	file, err := os.Open(filepath.Join(files...))
//	if err != nil {
//		return nil, err
//	}
//
//	r, err := invalid.NewRule(file)
//	if err != nil {
//		return nil, err
//	}
//
//	return r, nil
//
//}
//
//func TestLinting(t *testing.T) {
//	r, err := loadTestCase([]string{"yaml-cases", "simple.yaml"})
//	assert.Nil(t, err)
//	assert.NotNil(t, r)
//	r.Decode()
//
//	rule, err := loadTestRule([]string{"exam", "simple_exam2.yaml"})
//	assert.Nil(t, err)
//	assert.NotNil(t, r)
//	err = rule.Restructure()
//	assert.Nil(t, err)
//
//	rule.ValidYAML(r)
//
//}
