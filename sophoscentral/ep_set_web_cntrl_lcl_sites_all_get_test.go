package sophoscentral

import (
	"github.com/brianvoe/gofakeit/v6"
	"testing"
)

func TestWebControlLocalSiteUpdateRequest_Verify(t *testing.T) {
	faker := gofakeit.New(0)

	type fields struct {
		CategoryID int
		Tags       []string
		URL        string
		Comment    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				CategoryID: faker.Number(1, 57),
				Tags:       nil,
				URL:      	faker.URL(),
				Comment:    faker.LetterN(300),
			},
			wantErr: false,
		},
		{
			name: "No Tags, Category ID out of spec",
			fields: fields{
				CategoryID: 58,
				Tags:       nil,
				URL:      	faker.URL(),
				Comment:    faker.LetterN(300),
			},
			wantErr: true,
		},
		{
			name: "NoCategory ID, Tags out of spec",
			fields: fields{
				CategoryID: 58,
				Tags:       []string{faker.LetterN(51), faker.LetterN(150)},
				URL:      	faker.URL(),
				Comment:    faker.LetterN(300),
			},
			wantErr: true,
		},
		{
			name: "url too long",
			fields: fields{
				CategoryID: faker.Number(1, 57),
				Tags:       nil,
				URL:      	faker.LetterN(3000),
				Comment:    faker.LetterN(300),
			},
			wantErr: true,
		},
		{
			name: "no tags, no category",
			fields: fields{
				CategoryID: 0,
				Tags:       nil,
				URL:      	faker.LetterN(300),
				Comment:    faker.LetterN(300),
			},
			wantErr: true,
		},
		{
			name: "valid tags, no category",
			fields: fields{
				CategoryID: 0,
				Tags:       []string{faker.LetterN(1), faker.LetterN(50)},
				URL:      	faker.LetterN(300),
				Comment:    faker.LetterN(300),
			},
			wantErr: false,
		},
		{
			name: "comment too long",
			fields: fields{
				CategoryID: 0,
				Tags:       []string{faker.LetterN(1), faker.LetterN(50)},
				URL:      	faker.LetterN(300),
				Comment:    faker.LetterN(3000),
			},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WebControlLocalSiteRequest{
				CategoryID: tt.fields.CategoryID,
				Tags:       tt.fields.Tags,
				URL:        tt.fields.URL,
				Comment:    tt.fields.Comment,
			}
			if err := w.Verify(); (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
