package config

import (
	op "imagecore/operation" // Grab `EncoderOption` from operation package.
)

// This ts a utility function to convert pipeline block to image operation.

func PipelineBlockToOperation(pb PipelineBlock) op.Operation {

	// We assume the loaded config has been fully checked.
	// Therefore we don't return error here.
	switch pb.Operation {

	case "decode":
		return op.Decode()

		// - `icc_embed`

		// - `write`
	case "resize":
		if pb.Resize.Factor != 0.0 {
			return op.ResizeImageByFactor(pb.Resize.Algorithm, pb.Resize.Factor)
		} else if pb.Resize.Width != 0 {
			return op.ResizeImageByWidth(pb.Resize.Algorithm, pb.Resize.Width)
		} else if pb.Resize.Height != 0 {
			return op.ResizeImageByHeight(pb.Resize.Algorithm, pb.Resize.Height)
		} else {
			return nil // This should not happen, since the config has been checked.
		}

	case "crop":
		return op.Crop(pb.Crop.Width, pb.Crop.Height, pb.Crop.Alignment)

	case "encode":
		return op.Encode(pb.Encode.Format, (*op.EncoderOption)(pb.Encode.Options))

	case "icc_embed":
		return op.EmbedProfile(pb.ICCEmbedProfile.ProfileName)

	case "write":
		return op.WriteImageToFile(pb.Write.GenerateFileName())
	default:
		return nil // This should not happen, since the config has been checked.
	}
}